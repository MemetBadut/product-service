package model

import (
	"time"

	"gorm.io/gorm"
)

// Product adalah representasi tabel 'products' di database.
// GORM secara otomatis membuat tabel berdasarkan struct ini.
type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string         `gorm:"not null;size:255" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"not null" json:"price"`
	Stock       int            `gorm:"not null;default:0" json:"stock"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"` // Soft delete
}

func (Product) TableName() string {
	return "products"
}

type CreateProductRequest struct {
	Name 	  string  `json:"name" binding:"required"`
	Description string `json:"description"`
	Price 	  float64 `json:"price" binding:"required,gt=0"`
	Stock 	  int     `json:"stock" binding:"required,gte=0"`
}

type UpdateStockRequest struct {
	Quantity int `json:"quantity" validate:"required"`
}
