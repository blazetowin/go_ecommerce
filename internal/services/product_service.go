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
	product.InStock = product.Stock > 0
	return s.Repo.Create(product)
}


// Hatalı ürünler için özel hata tipi
type InvalidProductError struct{}

func (e *InvalidProductError) Error() string {
	return "Ürün adı ve fiyatı geçerli olmalıdır"
}

func (s *ProductService) GetFiltered(name string, minPrice, maxPrice, minStock, maxStock int) ([]models.Product, error) {
	return s.Repo.FindFilteredProducts(name, minPrice, maxPrice, minStock, maxStock)
}
