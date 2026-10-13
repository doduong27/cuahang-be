package handlers

import (
	"net/http"
	"os"
	"time"

	"cuahang-be/internal/database"

	"github.com/gin-gonic/gin"
)

var StartTime = time.Now()

// HealthCheck phục vụ Liveness và Readiness Probe cho Kubernetes
func HealthCheck(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		status := "healthy"
		dbStatus := "connected"

		if database.DB != nil {
			sqlDB, err := database.DB.DB()
			if err != nil || sqlDB.Ping() != nil {
				dbStatus = "disconnected"
				status = "degraded"
			}
		} else {
			dbStatus = "not_initialized"
			status = "degraded"
		}

		c.JSON(http.StatusOK, gin.H{
			"status":      status,
			"environment": env,
			"database":    dbStatus,
			"uptime":      time.Since(StartTime).String(),
			"timestamp":   time.Now().Format(time.RFC3339),
		})
	}
}

// AppInfo trả về thông tin runtime và Pod hostname
func AppInfo(env string) gin.HandlerFunc {
	return func(c *gin.Context) {
		hostname, _ := os.Hostname()
		c.JSON(http.StatusOK, gin.H{
			"app":         "cuahang-backend-go",
			"version":     "1.0.0",
			"runtime":     "Golang 1.23",
			"environment": env,
			"hostname":    hostname,
		})
	}
}
