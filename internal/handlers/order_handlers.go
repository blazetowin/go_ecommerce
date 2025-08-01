package handlers

import (
	"encoding/json"
	"net/http"

	"go_ecommerce/internal/services"
)

type OrderHandler struct {
	OrderService *services.OrderService
}

func NewOrderHandler(service *services.OrderService) *OrderHandler {
	return &OrderHandler{
		OrderService: service,
	}
}

func (h *OrderHandler) GetOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderService.GetAllOrders()
	if err != nil {
		http.Error(w, "Siparişler getirilemedi", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}

func (h *OrderHandler) GetOrderHistory(w http.ResponseWriter, r *http.Request) {
	orders, err := h.OrderService.GetOrderHistory()
	if err != nil {
		http.Error(w, "Sipariş geçmişi alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orders)
}
