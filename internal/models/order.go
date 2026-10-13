package models

import "time"

// Order định nghĩa đơn hàng
type Order struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	CustomerName string    `gorm:"not null" json:"customerName"`
	TotalAmount  float64   `gorm:"not null" json:"totalAmount"`
	Status       string    `gorm:"default:'confirmed'" json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
}

// OrderInput dữ liệu đầu vào khi đặt hàng
type OrderInput struct {
	CustomerName string  `json:"customerName" binding:"required"`
	TotalAmount  float64 `json:"totalAmount" binding:"required"`
}
