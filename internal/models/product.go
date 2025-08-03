package models

import "gorm.io/gorm"

type Product struct {
	gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt otomatik gelir
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Price       float64  `json:"price"`
	InStock     bool     `json:"in_stock" gorm:"-"`
	Stock 		int  	 `json:"stock"`
}
