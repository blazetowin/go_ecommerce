package repositories

import (
	"strings"

	"go_ecommerce/database"
	"go_ecommerce/internal/models"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		db: database.DB,
	}
}

// Tüm ürünleri getir
func (r *ProductRepository) GetAll() ([]models.Product, error) {
	var products []models.Product
	result := r.db.Find(&products)

	// Her ürünün InStock değerini hesapla
	for i := range products {
		products[i].InStock = products[i].Stock > 0
	}

	return products, result.Error
}

// Yeni ürün oluştur
func (r *ProductRepository) Create(product *models.Product) error {
	result := r.db.Create(product)
	return result.Error
}

// İsme göre stok getir
func (r *ProductRepository) GetStockByProductName(name string) int {
	var product models.Product
	err := r.db.Where("LOWER(name) = ?", strings.ToLower(name)).
		First(&product).Error
	if err != nil {
		return 0 // Hata varsa 0 dön
	}
	return product.Stock
}

// İsme göre stoğu güncelle
func (r *ProductRepository) UpdateStockByName(name string, newStock int) error {
	var product models.Product
	if err := r.db.Where("name = ?", name).First(&product).Error; err != nil {
		return err
	}
	product.Stock = newStock
	product.InStock = newStock > 0
	return r.db.Save(&product).Error
}

// ID ile stoğu güncelle
func (r *ProductRepository) UpdateStock(productID uint, newStock int) error {
	var product models.Product
	if err := r.db.First(&product, productID).Error; err != nil {
		return err
	}
	product.Stock = newStock
	product.InStock = newStock > 0
	return r.db.Save(&product).Error
}

// Filtreli ürün listeleme
func (r *ProductRepository) FindFilteredProducts(name string, minPrice, maxPrice, minStock, maxStock int) ([]models.Product, error) {
	db := r.db.Model(&models.Product{})

	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if minPrice > 0 {
		db = db.Where("price >= ?", minPrice)
	}
	if maxPrice > 0 {
		db = db.Where("price <= ?", maxPrice)
	}
	if minStock > 0 {
		db = db.Where("stock >= ?", minStock)
	}
	if maxStock > 0 {
		db = db.Where("stock <= ?", maxStock)
	}

	var products []models.Product
	if err := db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}
func (r *ProductRepository) DeleteByID(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *ProductRepository) GetByName(name string) (*models.Product, error) {
	var product models.Product
	result := r.db.Where("name = ?", name).First(&product)
	if result.Error != nil {
		return nil, result.Error
	}

	product.InStock = product.Stock > 0 // otomatik ayarla
	return &product, nil
}


func (r *ProductRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}
