package database

import (
	"log"

	"cuahang-be/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB(dsn string) (*gorm.DB, error) {
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[CẢNH BÁO] Không thể kết nối Database ngay lúc này: %v", err)
		return nil, err
	}

	log.Println("[Database] Kết nối PostgreSQL thành công!")

	// 1. Tự động Migration cấu hình bảng
	if err := DB.AutoMigrate(&models.Product{}, &models.Order{}); err != nil {
		log.Printf("[Database] Lỗi AutoMigrate: %v", err)
	}

	// 2. Tự động nạp dữ liệu mẫu nếu bảng Product đang trống
	var count int64
	DB.Model(&models.Product{}).Count(&count)
	if count == 0 {
		SeedProducts(DB)
	}

	return DB, nil
}
