package services

import (
	"go_ecommerce/internal/models"
	"go_ecommerce/internal/repositories"
)

type ProductService struct {
	Repo *repositories.ProductRepository
}

// Constructor
func NewProductService(repo *repositories.ProductRepository) *ProductService {
	return &ProductService{
		Repo: repo,
	}
}

// Tüm ürünleri getir
func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.Repo.GetAll()
}

// Yeni ürün oluştur
func (s *ProductService) CreateProduct(product *models.Product) error {
	// Basit bir validasyon örneği:
	if product.Name == "" || product.Price <= 0 {
		return &InvalidProductError{}
	}
	return s.Repo.Create(product)
}

// Hatalı ürünler için özel hata tipi
type InvalidProductError struct{}

func (e *InvalidProductError) Error() string {
	return "Ürün adı ve fiyatı geçerli olmalıdır"
}
