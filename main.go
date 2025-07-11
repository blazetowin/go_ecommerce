//Package clause
//Uygulamanın giriş noktası olan ana paket
package main

//Import statement
import (
    "go_ecommerce/router"
    "go_ecommerce/database"
)

//Kod bloğu
func main() {
    //router paketinden router'ı al
    r := router.SetupRouter()
    database.Connect()
    // HTTP sunucusunu 8080 portunda başlat
    // Tarayıcıdan http://localhost:8080 
    r.Run(":8080") 
}
