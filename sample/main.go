package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
)

func main() {
	r := gin.Default()
	sessions.InitFlash([]byte("abcdefghijklmn"))

	r.GET("/ping", func(c *gin.Context) {
		sessions.AddFlash(c,"hello world")
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/ping1", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run("0.0.0.0:9090") // listen and serve on 0.0.0.0:8080
}