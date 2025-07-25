package repositories

import (
	"go_ecommerce/internal/models"
	"go_ecommerce/database"
	"strings"
)

type ProductRepository struct{}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{}
}

// Tüm ürünleri getir
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	result := database.DB.Find(&products)
	return products, result.Error
}

// Yeni ürün oluştur
func (r *ProductRepository) Create(product *models.Product) error {
	result := database.DB.Create(product)
	return result.Error
}

//Kaç stok olduğuna bakalım
func GetStockByProductName(name string) int {
	var product models.Product
	err := database.DB.Where("LOWER(name) = ?", strings.ToLower(name)).
		First(&product).Error
	if err != nil {
		return 0 // Hata varsa 0 dön
	}
	return product.Stock
}
