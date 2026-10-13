package database

import (
	"log"

	"cuahang-be/internal/models"

	"gorm.io/gorm"
)

// SeedProducts chèn các sản phẩm mẫu vào PostgreSQL
func SeedProducts(db *gorm.DB) {
	sampleProducts := []models.Product{
		{Name: "Bàn Phím Cơ Không Dây RGB", Price: 850000, Stock: 15, Category: "Bàn phím", ImageURL: "https://images.unsplash.com/photo-1587829741301-dc798b83add3?w=500"},
		{Name: "Chuột Gaming Siêu Nhẹ 65g", Price: 550000, Stock: 20, Category: "Chuột", ImageURL: "https://images.unsplash.com/photo-1527864550417-7fd91fc51a46?w=500"},
		{Name: "Tai Nghe Chống Ồn ANC Pro", Price: 1200000, Stock: 8, Category: "Âm thanh", ImageURL: "https://images.unsplash.com/photo-1505740420928-5e560c06d30e?w=500"},
		{Name: "Màn Hình 27 inch 2K 165Hz", Price: 3900000, Stock: 5, Category: "Màn hình", ImageURL: "https://images.unsplash.com/photo-1527443224154-c4a3942d3acf?w=500"},
		{Name: "Giá Đỡ Laptop Nhôm Công Thái Học", Price: 250000, Stock: 30, Category: "Phụ kiện", ImageURL: "https://images.unsplash.com/photo-1588872657578-7efd1f1555ed?w=500"},
	}
	db.Create(&sampleProducts)
	log.Println("[Database] Đã khởi tạo 5 sản phẩm mẫu thành công!")
}
