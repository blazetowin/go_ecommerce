package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"regexp"
	"strings"
	"go_ecommerce/internal/models"
	"go_ecommerce/internal/repositories"
)

type ChatService struct {
	OrderRepo   *repositories.OrderRepository
	ProductRepo *repositories.ProductRepository
	APIKey      string
}

func NewChatService(orderRepo *repositories.OrderRepository, productRepo *repositories.ProductRepository, apiKey string) *ChatService {
	return &ChatService{
		OrderRepo:   orderRepo,
		ProductRepo: productRepo,
		APIKey:      apiKey,
	}
}

func (cs *ChatService) GetOrderRepo() *repositories.OrderRepository {
	return cs.OrderRepo
}

func (cs *ChatService) GetProductRepo() *repositories.ProductRepository {
	return cs.ProductRepo
}

// Kullanıcının mesajını yorumla ve cevap döndür
func (cs *ChatService) GetResponse(prompt string) string {
	prompt = strings.ToLower(prompt)

	// 1. Sipariş sorgusu öncelikli kontrol edilmeli
	if msg, ok := cs.CheckIfOrderHistoryQuery(prompt); ok {
		return msg
	}

	// 2. Satın alma niyeti varsa kontrol et
	if msg, ok := cs.CheckIfPurchaseIntent(prompt); ok {
		return msg
	}

	// 3. En son stok sorgusunu yap
	if msg, ok := cs.GetDynamicAnswer(prompt); ok {
		return msg
	}

	// 4. Hiçbiri değilse Gemini'ye sor
	reply, err := cs.AskQuestion(prompt)
	if err != nil {
		return "🤖 Bir hata oluştu: " + err.Error()
	}
	return reply
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

// Satın alma niyeti kontrolü
func (cs *ChatService) CheckIfPurchaseIntent(userInput string) (string, bool) {
	products, err := cs.ProductRepo.GetAll()
	if err != nil {
		return "Ürün bilgilerine ulaşılamıyor.", true
	}

	lowerInput := strings.ToLower(userInput)
	purchaseKeywords := []string{"satın almak", "satın al", "sipariş ver", "alacağım", "almak istiyorum", "siparişi ver"}

	for _, p := range products {
		if strings.Contains(lowerInput, strings.ToLower(p.Name)) {
			for _, keyword := range purchaseKeywords {
				if strings.Contains(lowerInput, keyword) {
					if p.Stock <= 0 {
						return fmt.Sprintf("Üzgünüz, şu anda %s stokta yok.", p.Name), true
					}

					err := cs.OrderRepo.CreateOrder(p.Name, 1)
					if err != nil {
						return "Sipariş oluşturulurken bir hata oluştu.", true
					}

					if err := cs.ProductRepo.UpdateStockByName(p.Name, p.Stock-1); err != nil {
						return "Stok güncellenemedi, siparişiniz alınamadı.", true
					}

					if p.Stock-1 <= 2 {
						return fmt.Sprintf("Siparişiniz başarıyla oluşturuldu! (%s)\n⚠️ Dikkat! Stokta yalnızca %d adet kaldı.", p.Name, p.Stock-1), true
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
	// Sadece açık bir şekilde stok soruluyorsa bu fonksiyon devreye girmeli
	if !(strings.Contains(prompt, "stokta var mı") ||
		strings.Contains(prompt, "stok durumu") ||
		strings.Contains(prompt, "kaç adet var") ||
		strings.Contains(prompt, "mevcut mu") ||
		strings.Contains(prompt, "var mı")) {
		return "", false
	}

	products, err := cs.ProductRepo.GetAll()
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

// "son X [ürün adı] sipariş" formatını yakalar
func ExtractLastNAndProduct(input string) (int, string) {
	re := regexp.MustCompile(`son (\d+)\s*([a-zA-Z0-9\s]*)?sipariş`)
	matches := re.FindStringSubmatch(strings.ToLower(input))

	n := 5 // default
	product := ""

	if len(matches) >= 2 && matches[1] != "" {
		if parsed, err := strconv.Atoi(matches[1]); err == nil {
			n = parsed
		}
	}
	if len(matches) >= 3 {
		product = strings.TrimSpace(matches[2])
	}
	return n, product
}

// Siparişleri kullanıcıya düzgün formatta göster
func FormatOrdersForDisplay(orders []models.Order) string {
	if len(orders) == 0 {
		return "📭 İstenilen kriterlere göre sipariş bulunamadı."
	}

	var sb strings.Builder
	sb.WriteString("🧾 Sipariş Geçmişiniz:\n")
	for i, o := range orders {
		sb.WriteString(fmt.Sprintf("%d. %s — %d adet — %s\n",
			i+1, o.ProductName, o.Quantity, o.CreatedAt.Format("2006-01-02 15:04")))
	}
	return sb.String()
}
func (cs *ChatService) CheckIfOrderHistoryQuery(prompt string) (string, bool) {
	if strings.Contains(prompt, "son") && strings.Contains(prompt, "sipariş") {
		n, product := ExtractLastNAndProduct(prompt)

		var orders []models.Order
		var err error

		if product != "" {
			orders, err = cs.OrderRepo.GetLastNOrdersByProduct(product, n)
		} else {
			orders, err = cs.OrderRepo.GetLastNOrders(n)
		}

		if err != nil {
			return "⚠️ Siparişler alınırken bir hata oluştu.", true
		}
		return FormatOrdersForDisplay(orders), true
	}
	return "", false
}

