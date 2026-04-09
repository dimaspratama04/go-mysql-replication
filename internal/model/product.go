package model

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"  json:"id"`
	Name        string         `gorm:"type:varchar(255);not null" json:"name"`
	Description string         `gorm:"type:text"                  json:"description"`
	Price       float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	Stock       int            `gorm:"default:0"                  json:"stock"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"             json:"created_at"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"             json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index"                      json:"deleted_at,omitempty"`
}

type CreateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type UpdateProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
