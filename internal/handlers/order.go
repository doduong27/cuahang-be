package handlers

import (
	"net/http"

	"cuahang-be/internal/database"
	"cuahang-be/internal/models"

	"github.com/gin-gonic/gin"
)

// CreateOrder lưu đơn hàng mới vào cơ sở dữ liệu
func CreateOrder(c *gin.Context) {
	if database.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database chưa sẵn sàng"})
		return
	}
	var input models.OrderInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dữ liệu đơn hàng không hợp lệ", "details": err.Error()})
		return
	}

	order := models.Order{
		CustomerName: input.CustomerName,
		TotalAmount:  input.TotalAmount,
		Status:       "confirmed",
	}
	if err := database.DB.Create(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi lưu đơn hàng vào PostgreSQL"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Đặt hàng thành công",
		"order":   order,
	})
}
