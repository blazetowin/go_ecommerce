package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"go_ecommerce/database"
	"go_ecommerce/internal/handlers"
	"go_ecommerce/internal/repositories"
	"go_ecommerce/internal/services"
	"go_ecommerce/utils"
)

func main() {
	utils.LoadEnv() // .env dosyasını yükle

	// 1. Veritabanı bağlantısı
	database.Connect()

	// 2. Repository → Service → Handler zincirini kur

	// 📦 Ürünler
	productRepo := repositories.NewProductRepository()
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// 🧾 Siparişler
	orderRepo := repositories.NewOrderRepository()
	orderService := services.NewOrderService()
	orderHandler := handlers.NewOrderHandler(orderService)

	// 🛒 Sepet
	cartRepo := repositories.NewCartRepository()
	// 💬 Chat (Sepet işlemleri de burada yapılacak)
	apiKey := os.Getenv("GEMINI_API_KEY")
	chatService := services.NewChatService(orderRepo, productRepo, apiKey)
	chatService.CartRepo = cartRepo // sepete erişim için repo ekle
	chatHandler := handlers.NewChatHandler(chatService)

	// 3. Router
	mux := http.NewServeMux()

	// 4. Routes

	// Ürün işlemleri
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productHandler.GetProducts(w, r)
		} else if r.Method == http.MethodPost {
			productHandler.CreateProduct(w, r)
		} else {
			http.Error(w, "Yalnızca GET ve POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// Ürün silme işlemi
	mux.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			// URL'den ID’yi al (örnek: /api/products/3)
			parts := strings.Split(r.URL.Path, "/")
			if len(parts) != 4 {
				http.Error(w, "Geçersiz istek yolu", http.StatusBadRequest)
				return
			}
			idStr := parts[3]
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "Geçersiz ID", http.StatusBadRequest)
				return
			}
			productHandler.DeleteProduct(w, r, uint(id))
		} else {
			http.Error(w, "Yalnızca DELETE destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 💬 Chat bot endpoint
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.HandleChat(w, r)
		} else {
			http.Error(w, "Yalnızca POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 🧾 Sipariş listeleme
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			orderHandler.GetOrders(w, r)
		} else {
			http.Error(w, "Yalnızca GET destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 🛒 Sepete ürün ekleme
	mux.HandleFunc("/api/cart/add", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.HandleAddToCart(w, r)
		} else {
			http.Error(w, "Yalnızca POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 🗑️ Sepetten ürün çıkarma
	mux.HandleFunc("/api/cart/remove", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.HandleRemoveFromCart(w, r)
		} else {
			http.Error(w, "Yalnızca POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 5. Server başlat
	port := ":8080"
	fmt.Println("🚀 Sunucu çalışıyor: http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}
