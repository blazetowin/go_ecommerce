package main

import (
	"fmt"
	"log"
	"net/http"

	"go_ecommerce/database"
	"go_ecommerce/internal/handlers"
	"go_ecommerce/internal/repositories"
	"go_ecommerce/internal/services"
	"go_ecommerce/utils"
)

func main() {
	utils.LoadEnv() // .env dosyasını yükle
	// Uygulama başlatma adımları:
	// 1. Veritabanı bağlantısını başlat
	database.Connect()

	// 2. Repository → Service → Handler zincirini kur
	productRepo := repositories.NewProductRepository()
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// 3. Chat servisi ve handler'ı kur
	chatService := services.NewChatService()
	chatHandler := handlers.NewChatHandler(chatService)

	// 4. Router oluştur
	mux := http.NewServeMux()

	// 5. Product endpoint
	mux.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			productHandler.GetAllProducts(w, r)
		} else if r.Method == http.MethodPost {
			productHandler.CreateProduct(w, r)
		} else {
			http.Error(w, "Yalnızca GET ve POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 6. Chat endpoint
	mux.HandleFunc("/api/chat", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			chatHandler.HandleChat(w, r)
		} else {
			http.Error(w, "Yalnızca POST destekleniyor", http.StatusMethodNotAllowed)
		}
	})

	// 7. Server'ı başlat
	port := ":8080"
	fmt.Println("🚀 Sunucu çalışıyor: http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, mux))
}
