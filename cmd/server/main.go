package main

import (
	"log"

	"cuahang-be/internal/config"
	"cuahang-be/internal/database"
	"cuahang-be/internal/handlers"
	"cuahang-be/internal/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.LoadConfig()

	// Khởi tạo kết nối Database và AutoMigrate
	_, _ = database.ConnectDB(cfg.DatabaseURL)

	if cfg.Environment == "prod" || cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	// 1. Healthcheck & Cluster Info
	r.GET("/healthz", handlers.HealthCheck(cfg.Environment))
	r.GET("/api/info", handlers.AppInfo(cfg.Environment))

	// 2. Business REST APIs
	r.GET("/api/products", handlers.GetProducts)
	r.GET("/api/products/:id", handlers.GetProductByID)
	r.POST("/api/orders", handlers.CreateOrder)

	log.Printf("[cuahang-be] Go Backend đang chạy tại http://0.0.0.0:%s (Môi trường: %s)", cfg.Port, cfg.Environment)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Không thể khởi động server: %v", err)
	}
}
