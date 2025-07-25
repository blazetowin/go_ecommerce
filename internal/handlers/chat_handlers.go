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

	// ✨ Önce veritabanına bakarak cevabı üret
	dynamicAnswer, matched := services.GetDynamicAnswer(userInput)
	if matched {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ChatResponse{Answer: dynamicAnswer})
		return
	}

	// 🧠 Eğer eşleşme yoksa Gemini'den cevap al
	answer, err := h.ChatService.AskQuestion(userInput)
	if err != nil {
		fmt.Println("❌ OpenAI ile konuşma hatası:", err)
		http.Error(w, "Gemini yanıtı alınamadı", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ChatResponse{Answer: answer})
}
