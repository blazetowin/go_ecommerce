package models

import (
	"gorm.io/gorm"
	"time"
)

type Cart struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint
	ProductID uint
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
