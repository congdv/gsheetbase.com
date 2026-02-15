package main

import (
	"log"
	"net/http"
	"time"

	"gsheetbase/shared/database"
	"gsheetbase/shared/repository"
	"gsheetbase/web/internal/config"
	"gsheetbase/web/internal/http/handlers"
	"gsheetbase/web/internal/http/middleware"
	"gsheetbase/web/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	db, err := database.Connect(cfg.DBURL)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}
	defer db.Close()

	// Repositories
	userRepo := repository.NewUserRepo(db)
	allowedSheetRepo := repository.NewAllowedSheetRepo(db)
	usageRepo := repository.NewUsageRepo(db)

	// Services
	authService := services.NewAuthService(cfg, userRepo)
	sheetService := services.NewSheetService(cfg, allowedSheetRepo)

	r := gin.Default()

	// Serve static frontend files
	r.Static("/assets", "./web/ui/dist/assets")
	r.StaticFile("/favicon.ico", "./web/ui/dist/favicon.ico")
	r.StaticFile("/robots.txt", "./web/ui/dist/robots.txt")
	r.StaticFile("/site.webmanifest", "./web/ui/dist/site.webmanifest")
	r.StaticFile("/sitemap.xml", "./web/ui/dist/sitemap.xml")

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.FrontendOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type", "X-API-Key"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "ok")
	})

	// Public config endpoint for frontend
	configHandler := handlers.NewFrontendConfigHandler(cfg)
	r.GET("/frontend-config", configHandler.GetConfig)

	api := r.Group("/api")

	// Google OAuth routes
	googleHandler := handlers.NewGoogleAuthHandler(authService, cfg)
	api.GET("/auth/google/start", googleHandler.Start)
	api.GET("/auth/google/callback", googleHandler.Callback)
	api.POST("/auth/google/request-scopes", middleware.Authenticate(cfg, authService), googleHandler.RequestAdditionalScopes)

	// Authenticated routes (JWT from Google OAuth)
	authHandler := handlers.NewAuthHandler(authService, cfg)
	api.GET("/auth/me", middleware.Authenticate(cfg, authService), authHandler.Me)
	api.POST("/auth/logout", authHandler.Logout)

	// Session refresh (uses refresh_token cookie, no auth middleware needed)
	refreshAuthHandler := handlers.NewRefreshAuthHandler(authService, cfg)
	api.POST("/auth/refresh-session", refreshAuthHandler.RefreshSession)

	// Sheet registration (must register sheets before accessing them)
	allowedSheetHandler := handlers.NewAllowedSheetHandler(allowedSheetRepo)
	api.POST("/sheets/register", middleware.Authenticate(cfg, authService), allowedSheetHandler.Register)
	api.GET("/sheets/registered", middleware.Authenticate(cfg, authService), allowedSheetHandler.List)
	api.DELETE("/sheets/registered/:sheet_id", middleware.Authenticate(cfg, authService), allowedSheetHandler.Delete)
	api.POST("/sheets/:id/publish", middleware.Authenticate(cfg, authService), allowedSheetHandler.Publish)
	api.DELETE("/sheets/:id/unpublish", middleware.Authenticate(cfg, authService), allowedSheetHandler.Unpublish)
	api.PATCH("/sheets/:id/write-settings", middleware.Authenticate(cfg, authService), allowedSheetHandler.UpdateWriteSettings)

	// Authentication management (bearer token and basic auth setup)
	api.GET("/sheets/:id/auth", middleware.Authenticate(cfg, authService), allowedSheetHandler.GetAuthStatus)
	api.POST("/sheets/:id/auth/type", middleware.Authenticate(cfg, authService), allowedSheetHandler.SetAuthType)
	api.POST("/sheets/:id/auth/bearer", middleware.Authenticate(cfg, authService), allowedSheetHandler.GenerateBearerToken)
	api.POST("/sheets/:id/auth/bearer/rotate", middleware.Authenticate(cfg, authService), allowedSheetHandler.RotateBearerToken)
	api.POST("/sheets/:id/auth/basic", middleware.Authenticate(cfg, authService), allowedSheetHandler.SetBasicAuth)
	api.DELETE("/sheets/:id/auth", middleware.Authenticate(cfg, authService), allowedSheetHandler.DisableAuth)

	// Sheet access (requires JWT auth + sheet must be registered)
	sheetHandler := handlers.NewSheetHandler(sheetService)
	api.POST("/sheets/create", middleware.Authenticate(cfg, authService), sheetHandler.CreateSheet)
	api.POST("/sheets/data", middleware.Authenticate(cfg, authService), sheetHandler.Get)

	// Analytics endpoints
	analyticsHandler := handlers.NewAnalyticsHandler(usageRepo, allowedSheetRepo)
	api.GET("/sheets/:id/analytics", middleware.Authenticate(cfg, authService), analyticsHandler.GetSheetAnalytics)
	api.GET("/analytics", middleware.Authenticate(cfg, authService), analyticsHandler.GetUserAnalytics)

	// Subscription & billing endpoints
	subscriptionHandler := handlers.NewSubscriptionHandler(userRepo, usageRepo)
	api.GET("/subscription/plan", middleware.Authenticate(cfg, authService), subscriptionHandler.GetCurrentPlan)
	api.GET("/subscription/usage", middleware.Authenticate(cfg, authService), subscriptionHandler.GetCurrentUsage)
	api.GET("/subscription/plans", subscriptionHandler.GetAvailablePlans) // Public endpoint

	// Serve index.html for all other routes (SPA fallback)
	r.NoRoute(func(ctx *gin.Context) {
		ctx.File("./web/ui/dist/index.html")
	})

	addr := ":" + cfg.Port
	log.Printf("API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
