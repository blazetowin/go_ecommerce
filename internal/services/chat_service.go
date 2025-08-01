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
	orderRepo   *repositories.OrderRepository
	productRepo *repositories.ProductRepository
	APIKey      string
}

func NewChatService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository, apiKey string) *ChatService {
	return &ChatService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		APIKey:      apiKey,
	}
}

// Gemini API'ye prompt gönder
func (cs *ChatService) AskQuestion(prompt string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1/models/gemini-1.5-flash:generateContent?key=" + cs.APIKey

	reqBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"role": "user",
				"parts": []map[string]string{
					{"text": prompt},
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

// Kullanıcının satın alma niyetini kontrol et
func (cs *ChatService) CheckIfPurchaseIntent(userInput string) (string, bool) {
	products, err := cs.productRepo.GetAll()
	if err != nil {
		return "Ürün bilgilerine ulaşılamıyor.", true
	}

	lowerInput := strings.ToLower(userInput)

	// Satın alma niyetini yansıtan kalıplar
	purchaseKeywords := []string{
		"satın almak istiyorum",
		"satın al",
		"almak istiyorum",
		"sipariş ver",
		"siparişi ver",
		"sipariş etmek istiyorum",
		"satın alma",
	}

	for _, p := range products {
		if strings.Contains(lowerInput, strings.ToLower(p.Name)) {
			for _, keyword := range purchaseKeywords {
				if strings.Contains(lowerInput, keyword) {
					if p.Stock <= 0 {
						return fmt.Sprintf("Üzgünüz, şu anda %s stokta yok.", p.Name), true
					}

					if err := cs.orderRepo.CreateOrder(p.Name, 1); err != nil {
						return "Sipariş oluşturulurken bir hata oluştu.", true
					}

					if err := cs.productRepo.UpdateStockByName(p.Name, p.Stock-1); err != nil {
						return "Stok güncellenemedi, siparişiniz alınamadı.", true
					}

					return fmt.Sprintf("Siparişiniz başarıyla oluşturuldu! (%s)", p.Name), true
				}
			}
		}
	}

	return "", false
}


// Dinamik stok sorgusu
func (cs *ChatService) GetDynamicAnswer(prompt string) (string, bool) {
	products, err := cs.productRepo.GetAll()
	if err != nil {
		return "Ürün bilgilerine ulaşılamadı.", true
	}

	for _, product := range products {
		if strings.Contains(strings.ToLower(prompt), strings.ToLower(product.Name)) {
			if product.Stock > 0 {
				return fmt.Sprintf("Evet, stokta %d adet %s var.", product.Stock, product.Name), true
			}
			return fmt.Sprintf("Üzgünüz, şu anda %s stokta yok.", product.Name), true
		}
	}

	return "", false
}
