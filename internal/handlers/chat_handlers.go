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

	userInput := req.Prompt

	// 🛒 Önce satın alma niyeti kontrol edilir
	if purchaseResponse, matched := h.ChatService.CheckIfPurchaseIntent(userInput); matched {
		json.NewEncoder(w).Encode(ChatResponse{Answer: purchaseResponse})
		return
	}

	// 📦 Sonra stok bilgisi kontrol edilir
	if dynamicAnswer, matched := h.ChatService.GetDynamicAnswer(userInput); matched {
		json.NewEncoder(w).Encode(ChatResponse{Answer: dynamicAnswer})
		return
	}

	// 🤖 Son olarak Gemini'den genel cevap alınır
	answer, err := h.ChatService.AskQuestion(userInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(ChatResponse{Answer: answer})
}
