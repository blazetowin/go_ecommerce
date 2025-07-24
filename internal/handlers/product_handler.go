package handlers

import (
	"encoding/json"
	"net/http"

	"go_ecommerce/internal/models"
	"go_ecommerce/internal/services"
)

type ProductHandler struct {
	ProductService *services.ProductService
}

// Constructor
func NewProductHandler(service *services.ProductService) *ProductHandler {
	return &ProductHandler{
		ProductService: service,
	}
}

// GET /api/products
func (h *ProductHandler) GetAllProducts(w http.ResponseWriter, r *http.Request) {
	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		http.Error(w, "Ürünler alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// POST /api/products
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	if err := h.ProductService.CreateProduct(&product); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(product)
}
