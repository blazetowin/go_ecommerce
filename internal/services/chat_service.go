package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"go_ecommerce/internal/repositories"
)

type ChatService struct {
	orderRepo *repositories.OrderRepository
	APIKey string
}


func NewChatService(orderRepo *repositories.OrderRepository, apiKey string) *ChatService {
	return &ChatService{
		orderRepo: orderRepo,
		APIKey:    apiKey, // ✔️ Artık parametre olarak geldiği için tanımlı
	}
}



func (cs *ChatService) AskQuestion(prompt string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key="+ cs.APIKey

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{
						"text": prompt,
					},
				},
			},
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("JSON oluşturulamadı: %v", err)
	}

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("İstek hatası: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("❌ Gemini API Hatası: %s", string(body))
	}

	var data struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("Cevap parse edilemedi: %v", err)
	}

	if len(data.Candidates) == 0 || len(data.Candidates[0].Content.Parts) == 0 {
		return "", fmt.Errorf("Gemini boş döndü: %s", string(body))
	}

	return data.Candidates[0].Content.Parts[0].Text, nil
}

func GetDynamicAnswer(userInput string)(string,bool){
	//Bu cümlede ürün adı geçiyor mu kontrol ediyor
	if strings.Contains(strings.ToLower(userInput), "iphone 14") && strings.Contains(userInput, "stok") {
		adet := repositories.GetStockByProductName("iPhone 14")
		if adet > 0 {
			return "Evet, stokta " + fmt.Sprintf("%d", adet) + " adet iPhone 14 var.", true
		}
		return "Maalesef şu anda iPhone 14 stokta yok.", true
	}
	return "", false
}

var orderRepo = repositories.NewOrderRepository()

func (s *ChatService) CheckIfPurchaseIntent(userInput string) (string, bool) {
	if strings.Contains(strings.ToLower(userInput), "satın almak istiyorum") && strings.Contains(userInput, "iphone 14") {
		currentStock := repositories.GetStockByProductName("iPhone 14")
		if currentStock <= 0 {
			return "Üzgünüz, şu anda iPhone 14 stokta yok.", true
		}

		err := s.orderRepo.CreateOrder("iPhone 14", 1)
		if err != nil {
			return "Sipariş oluşturulurken bir hata oluştu.", true
		}

		repositories.UpdateStock("iPhone 14", currentStock-1)

		return "Siparişiniz başarıyla oluşturuldu! 📦", true
	}
	return "", false
}

