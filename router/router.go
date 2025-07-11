// Bu dosya router adlı bir pakete ait
package router

import (
    "github.com/gin-gonic/gin" // Gin web framework'ünü içe aktar
)

// SetupRouter fonksiyonu *gin.Engine döner, tüm HTTP rotalarını burada tanımlarız
func SetupRouter() *gin.Engine {
    r := gin.Default() // Logger ve Recovery middlewareleri otomatik tanımlı bir engine döner



 // GET / isteğine yanıt veren basit bir test rotası
    r.GET("/", func(c *gin.Context) {

        // JSON tipinde 200 HTTP cevabı döner yani olumlu
        c.JSON(200, gin.H{
            "message": "E-Ticaret API'ye hoş geldin!",
        })
    })

    return r // Ayarlanmış router objesini geri döner
}
