//Package clause
//Uygulamanın giriş noktası olan ana paket
package main

//Import statement
import (
    "e_commerce/router"
)

func main() {
    //router paketinden router'ı al
    r := router.SetupRouter()

    // HTTP sunucusunu 8080 portunda başlat
    // Tarayıcıdan http://localhost:8080 
    r.Run(":8080") 
}
