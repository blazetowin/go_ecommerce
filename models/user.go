package models // Models paketi: Veritabanı modellerinin yer aldığı katman

import "gorm.io/gorm" z// GORM kütüphanesinin temel özelliklerini içe aktarır


// User yapısı: "users" tablosunu temsil eder
type User struct {
    gorm.Model 
	// gorm.Model: GORM tarafından otomatik eklenen alanlar:
    // - ID (primary key)
    // - CreatedAt (oluşturulma tarihi)
    // - UpdatedAt (güncellenme tarihi)
    // - DeletedAt (soft delete için, veriyi silmez; silinmiş gibi işaretler)
    Name     string `json:"name"`
    Email    string `json:"email" gorm:"unique"`//Aynı mail iki kez kaydedilmesin
    Password string `json:"password"`//Şifrelenmiş şekilde saklanmalı!
}
