package main

import (
	"log"
	"net/http"
	"time"

	"gsheetbase/shared/database"
	"gsheetbase/shared/repository"
	"gsheetbase/worker/internal/config"
	"gsheetbase/worker/internal/http/handlers"
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

	// Handlers
	sheetHandler := handlers.NewSheetHandler(sheetRepo, userRepo)

	// Setup Gin
	r := gin.Default()

	// CORS - permissive for public API
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// Public API routes
	v1 := r.Group("/v1")
	v1.GET("/:api_key", sheetHandler.GetPublic)

	addr := ":" + cfg.Port
	log.Printf("Worker API listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
