package main

import (
	"log"
	"net/http"
	"time"

	"gsheetbase/shared/database"
	"gsheetbase/shared/repository"
	"gsheetbase/worker/internal/cache"
	"gsheetbase/worker/internal/config"
	"gsheetbase/worker/internal/http/handlers"
	"gsheetbase/worker/internal/middleware"
	"gsheetbase/worker/internal/services"

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
	sheetRepo := repository.NewAllowedSheetRepo(db)
	userRepo := repository.NewUserRepo(db)
	usageRepo := repository.NewUsageRepo(db)

	// Optional: Redis and rate limiting (only if REDIS_URL is set)
	var rateLimitService *services.RateLimitService
	if cfg.RedisURL != "" {
		redisClient, err := cache.NewRedisClient(cfg.RedisURL)
		if err != nil {
			log.Printf("WARNING: Redis connection failed, rate limiting disabled: %v", err)
		} else {
			defer redisClient.Close()
			log.Println("Redis connected - rate limiting enabled")
			rateLimitService = services.NewRateLimitService(
				redisClient.GetClient(),
			)
		}
	} else {
		log.Println("Redis URL not configured - rate limiting disabled")
	}

	// Usage tracker with background workers
	usageTracker := middleware.NewUsageTracker(usageRepo, cfg.UsageTrackWorkers)
	defer usageTracker.Shutdown()

	// Handlers
	sheetHandler := handlers.NewSheetHandler(sheetRepo, userRepo)

	// Setup Gin
	r := gin.Default()

	// CORS - permissive for public API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public API routes with quota enforcement (rate limits + daily/monthly quotas)
	v1 := r.Group("/v1")

	// Group all /v1/:api_key routes and apply all middlewares only to them
	authService := services.NewAuthService()
	v1ApiKey := v1.Group(":api_key")
	if rateLimitService != nil {
		v1ApiKey.Use(middleware.QuotaEnforcementMiddleware(rateLimitService, usageRepo, userRepo, sheetRepo))
	}
	v1ApiKey.Use(middleware.UsageTrackingMiddleware(usageTracker))
	v1ApiKey.Use(middleware.AccessTokenEnsureMiddleware(sheetRepo, userRepo, authService, cfg.GoogleClientID, cfg.GoogleClientSecret))

	v1ApiKey.GET("", sheetHandler.GetPublic)
	v1ApiKey.POST("", sheetHandler.PostPublic)
	v1ApiKey.PUT("", sheetHandler.PutPublic)
	v1ApiKey.PATCH("", sheetHandler.PatchPublic)
	v1ApiKey.DELETE("", sheetHandler.DeletePublic)

	addr := ":" + cfg.Port
	log.Printf("Worker API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
