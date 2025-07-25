package database

import(
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
		log.Fatal("Veritabanina bağlanilamadi: ", err)
	}

	// GORM ile tabloları oluştur
	DB.AutoMigrate(&models.Product{})

	// 📦 Varsayılan ürünleri yükle
	var count int64
	DB.Model(&models.Product{}).Count(&count)
	if count == 0 {
		DB.Create(&models.Product{
			Name:  "iPhone 14",
			Stock: 5,
		})
	}
}
