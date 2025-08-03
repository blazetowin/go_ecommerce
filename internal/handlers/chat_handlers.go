package handlers

import (
	"encoding/json"
	"net/http"
	"go_ecommerce/internal/services"
)

type ChatHandler struct {
	ChatService *services.ChatService
}
func NewChatHandler(service *services.ChatService) *ChatHandler {
	return &ChatHandler{
		ChatService: service,
	}
}
type ChatRequest struct {
	Prompt string `json:"prompt"`
	UserID uint   `json:"user_id"`
}
type ChatResponse struct {
	Answer string `json:"answer"`
}
func (h *ChatHandler) HandleChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	response := h.ChatService.GetResponse(req.Prompt, req.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Answer: response})
}
type CartActionRequest struct {
	UserID    uint `json:"user_id"`
	ProductID uint `json:"product_id"`
	Quantity  int  `json:"quantity"` // Remove için opsiyonel
}

func (h *ChatHandler) HandleAddToCart(w http.ResponseWriter, r *http.Request) {
	var req CartActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	message, err := h.ChatService.AddToCart(req.UserID, req.ProductID, req.Quantity)
	if err != nil {
		http.Error(w, message, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": message})
}


func (h *ChatHandler) HandleRemoveFromCart(w http.ResponseWriter, r *http.Request) {
	var req CartActionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Geçersiz istek", http.StatusBadRequest)
		return
	}

	if err := h.ChatService.RemoveFromCart(req.UserID, req.ProductID); err != nil {
		http.Error(w, "Sepetten kaldırılamadı", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "🗑️ Sepetten çıkarıldı"})
}

