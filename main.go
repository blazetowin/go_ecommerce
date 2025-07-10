package main

import (
    "e_commerce/router"
)

func main() {
    r := router.SetupRouter()
    r.Run(":8080") // localhost:8080 portunda çalışır
}
