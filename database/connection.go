package database// Database paketi: Veritabanı bağlantı işlemlerini içerir

import (
    "ecommerce/models" // Veritabanı modellerini içe aktar (örneğin: User)
    "gorm.io/driver/sqlite" // SQLite sürücüsü (veritabanı motoru)
    "gorm.io/gorm" // GORM ORM kütüphanesi
    "log" // Hataları yazdırmak için kullanılır
)

var DB *gorm.DB // DB: Global olarak erişilebilecek veritabanı bağlantı nesnesi

func Connect() { // Connect fonksiyonu: Veritabanına bağlanır ve gerekli tabloları oluşturur

	// GORM ile SQLite veritabanına bağlan (dosya olarak çalışır)
    db, err := gorm.Open(sqlite.Open("ecommerce.db"), &gorm.Config{})
    if err != nil {
        log.Fatal("Veritabanina bağlanilamadi: ", err)
		 // Bağlantı başarısızsa uygulamayı durdur
    }

	// Veritabanı şemasıyla model eşleştirmesi yapılır
    // Eğer "users" tablosu yoksa otomatik olarak oluşturulur
    db.AutoMigrate(&models.User{})

	// Bağlantı başarılıysa global DB değişkenine atılır
    DB = db
}
