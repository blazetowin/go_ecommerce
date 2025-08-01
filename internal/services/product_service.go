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
// internal/services/product_service.go

func (s *ProductService) DeleteProductByID(id uint) error {
	return s.Repo.DeleteByID(id)
}

func (s *ProductService) UpsertProduct(newProduct *models.Product) (*models.Product, error) {
	existingProduct, err := s.Repo.GetByName(newProduct.Name)
	if err == nil && existingProduct != nil {
		// Ürün varsa: stoğunu güncelle
		existingProduct.Stock += newProduct.Stock
		if err := s.Repo.Update(existingProduct); err != nil {
			return nil, err
		}
		return existingProduct, nil
	}

	// Ürün yoksa: yeni ürün oluştur
	err = s.Repo.Create(newProduct)
	if err != nil {
		return nil, err
	}
	return newProduct, nil
}


