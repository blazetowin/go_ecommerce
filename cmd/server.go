package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
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

	// 💬 Chat
	orderRepo :=repositories.NewOrderRepository()
	apiKey := os.Getenv("GEMINI_API_KEY")
	chatService := services.NewChatService(orderRepo,apiKey)
	chatHandler := handlers.NewChatHandler(chatService)

	// 🧾 Siparişler
	orderService := services.NewOrderService()
	orderHandler := handlers.NewOrderHandler(orderService)

	// 3. Router
	mux := http.NewServeMux()

	// 4. Routes

	// 🚚 Ürün işlemleri
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productHandler.GetAllProducts(w, r)
		} else if r.Method == http.MethodPost {
			productHandler.CreateProduct(w, r)
		} else {
			http.Error(w, "Yalnızca GET ve POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 🤖 Chat bot endpoint
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.HandleChat(w, r)
		} else {
			http.Error(w, "Yalnızca POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 🧾 Sipariş listeleme endpointi
	mux.HandleFunc("/api/orders", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			orderHandler.GetOrders(w, r)
		} else {
			http.Error(w, "Yalnızca GET destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 5. Server başlat
	port := ":8080"
	fmt.Println("🚀 Sunucu çalışıyor: http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}
