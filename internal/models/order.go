package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	ProductName string `json:"product_name"`
	Quantity    int    `json:"quantity"`
	UserID      uint   `json:"user_id"`
}
