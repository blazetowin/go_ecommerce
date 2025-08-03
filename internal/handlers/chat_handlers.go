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

	// 🌟 Artık tüm senaryolar ChatService içinde yönetiliyor
	response := h.ChatService.GetResponse(req.Prompt, req.UserID)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Answer: response})
}
