package handlers

import (
	"net/http"
	"strconv"

	"cuahang-be/internal/database"
	"cuahang-be/internal/models"

	"github.com/gin-gonic/gin"
)

// GetProducts lấy toàn bộ danh sách sản phẩm
func GetProducts(c *gin.Context) {
	if database.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database chưa sẵn sàng"})
		return
	}
	var products []models.Product
	if err := database.DB.Order("id asc").Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Lỗi truy vấn dữ liệu"})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProductByID lấy thông tin 1 sản phẩm theo ID
func GetProductByID(c *gin.Context) {
	if database.DB == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Database chưa sẵn sàng"})
		return
	}
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID không hợp lệ"})
		return
	}
	var product models.Product
	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Không tìm thấy sản phẩm"})
		return
	}
	c.JSON(http.StatusOK, product)
}
