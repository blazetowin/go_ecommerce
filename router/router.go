package router

import (
    "github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
    r := gin.Default()

    // Basit test rotası
    r.GET("/", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "E-Ticaret API'ye hoş geldin!",
        })
    })

    return r
}
