package repositories

import (
	"fmt"
	"strings"
	"go_ecommerce/database"
	"go_ecommerce/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		db: database.DB,
	}
}

func (r *OrderRepository) CreateOrder(productName string, quantity int) error {
	var product models.Product
	if err := r.db.Where("name = ?", productName).First(&product).Error; err != nil {
		return err
	}

	if product.Stock < quantity {
		return fmt.Errorf("Stok yetersiz")
	}

	order := models.Order{
		ProductName: productName,
		Quantity:    quantity,
	}

	if err := r.db.Create(&order).Error; err != nil {
		return err
	}

	product.Stock -= quantity
	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Order("created_at desc").Find(&orders).Error
	return orders, err
}

func (r *OrderRepository) GetLastNOrders(n int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Order("created_at desc").Limit(n).Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (r *OrderRepository) GetLastNOrdersByProduct(product string, n int) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("LOWER(product_name) LIKE ?", "%"+strings.ToLower(product)+"%").
		Order("created_at desc").
		Limit(n).
		Find(&orders).Error
	if err != nil {
		return nil, err
	}
	return orders, nil
}
