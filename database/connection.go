package database

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"go_ecommerce/internal/models"
)

var DB *gorm.DB

func Connect() {
	var err error

	DB, err = gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ Veritabanı bağlantısı kurulamadı:", err)
	}

	fmt.Println("✅ Veritabanı bağlantısı başarılı")

	// Otomatik tablo oluştur
	err = DB.AutoMigrate(&models.Product{})
	if err != nil {
		log.Fatal("❌ Tablo migrate edilemedi:", err)
	}
}
