package handlers

import (
	"encoding/json"
	"fmt"
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
		http.Error(w, "Geçersiz veri", http.StatusBadRequest)
		return
	}

	userInput := req.Prompt

	// 🛒 1. Sipariş isteği kontrolü
	if purchaseAnswer, matched := h.ChatService.CheckIfPurchaseIntent(userInput); matched {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Answer: purchaseAnswer})
		return
	}

	// 📦 2. Stok kontrolü
	if dynamicAnswer, matched := h.ChatService.GetDynamicAnswer(userInput); matched {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Answer: dynamicAnswer})
		return
	}

	// 🤖 3. AI yanıtı (Gemini)
	answer, err := h.ChatService.AskQuestion(userInput)
	if err != nil {
		fmt.Println("❌ Gemini API hatası:", err)
		http.Error(w, "Gemini yanıtı alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Answer: answer})
}
