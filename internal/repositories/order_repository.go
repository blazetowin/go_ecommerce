package repositories

import (
	"fmt"
	"go_ecommerce/database"
	"go_ecommerce/internal/models"
	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository() *OrderRepository {
	return &OrderRepository{
		db: database.DB, // 🔗 db bağlantısını içeri alıyoruz
	}
}

func (r *OrderRepository) CreateOrder(productName string, quantity int) error {
	var product models.Product
	if err := r.db.Where("name = ?", productName).First(&product).Error; err != nil {
		return err
	}

	// ❗ Yeterli stok yoksa hata dön
	if product.Stock < quantity {
		return fmt.Errorf("Stok yetersiz")
	}

	// ✅ Siparişi oluştur
	order := models.Order{
		ProductName: productName,
		Quantity:    quantity,
	}

	if err := r.db.Create(&order).Error; err != nil {
		return err
	}

	// 🔻 Stoğu azalt
	product.Stock -= quantity
	if err := r.db.Save(&product).Error; err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetAllOrders() ([]models.Order, error) {
	var orders []models.Order
	result := r.db.Find(&orders)
	return orders, result.Error
}
