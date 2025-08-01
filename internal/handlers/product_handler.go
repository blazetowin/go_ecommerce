package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

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

// 🔁 GET /api/products → tüm ürünleri getir veya filtre uygula
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	// Filtre parametrelerini oku
	name := q.Get("name")
	minPrice, _ := strconv.Atoi(q.Get("min_price"))
	maxPrice, _ := strconv.Atoi(q.Get("max_price"))
	minStock, _ := strconv.Atoi(q.Get("min_stock"))
	maxStock, _ := strconv.Atoi(q.Get("max_stock"))

	// Eğer filtre parametresi varsa filtreli getir
	if name != "" || minPrice > 0 || maxPrice > 0 || minStock > 0 || maxStock > 0 {
		products, err := h.ProductService.GetFiltered(name, minPrice, maxPrice, minStock, maxStock)
		if err != nil {
			http.Error(w, "Ürünler getirilemedi", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
		return
	}

	// Aksi halde tüm ürünleri getir
	products, err := h.ProductService.GetAllProducts()
	if err != nil {
		http.Error(w, "Ürünler alınamadı", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

// 🆕 POST /api/products → ürün ekle
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	createdProduct, err := h.ProductService.UpsertProduct(&product)
	if err != nil {
		http.Error(w, "Ürün kaydedilemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdProduct)
}
// 🆔 GET /api/products/{id} → ID ile ürün getir
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request, id uint) {
	err := h.ProductService.DeleteProductByID(id)
	if err != nil {
		http.Error(w, "Ürün silinemedi: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ürün başarıyla silindi"))
}

