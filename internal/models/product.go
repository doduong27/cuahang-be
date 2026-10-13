package models

import "time"

// Product định nghĩa thông tin sản phẩm
type Product struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Price     float64   `gorm:"not null" json:"price"`
	Stock     int       `gorm:"default:10" json:"stock"`
	Category  string    `gorm:"default:'Phụ kiện'" json:"category"`
	ImageURL  string    `json:"imageUrl"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
