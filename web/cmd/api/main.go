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

	// Services
	authService := services.NewAuthService(cfg, userRepo)
	sheetService := services.NewSheetService(cfg, allowedSheetRepo)

	r := gin.Default()

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

	api := r.Group("/api")

	// Google OAuth routes
	googleHandler := handlers.NewGoogleAuthHandler(authService, cfg)
	api.GET("/auth/google/start", googleHandler.Start)
	api.GET("/auth/google/callback", googleHandler.Callback)

	// Authenticated routes (JWT from Google OAuth)
	authHandler := handlers.NewAuthHandler(authService, cfg)
	api.GET("/auth/me", middleware.Authenticate(cfg, authService), authHandler.Me)
	api.POST("/auth/logout", authHandler.Logout)

	// Sheet registration (must register sheets before accessing them)
	allowedSheetHandler := handlers.NewAllowedSheetHandler(allowedSheetRepo)
	api.POST("/sheets/register", middleware.Authenticate(cfg, authService), allowedSheetHandler.Register)
	api.GET("/sheets/registered", middleware.Authenticate(cfg, authService), allowedSheetHandler.List)
	api.DELETE("/sheets/registered/:sheet_id", middleware.Authenticate(cfg, authService), allowedSheetHandler.Delete)
	api.POST("/sheets/:id/publish", middleware.Authenticate(cfg, authService), allowedSheetHandler.Publish)
	api.DELETE("/sheets/:id/unpublish", middleware.Authenticate(cfg, authService), allowedSheetHandler.Unpublish)

	// Sheet access (requires JWT auth + sheet must be registered)
	sheetHandler := handlers.NewSheetHandler(sheetService)
	api.POST("/sheets/data", middleware.Authenticate(cfg, authService), sheetHandler.Get)

	addr := ":" + cfg.Port
	log.Printf("API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
