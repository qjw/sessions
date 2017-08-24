package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"net/http"
)

func main() {
	r := gin.Default()
	sessions.InitFlash([]byte("abcdefghijklmn"))

	r.GET("/ping", func(c *gin.Context) {
		sessions.AddFlash(c, "hello world")
		c.Redirect(http.StatusFound, "/pong")
	})

	r.GET("/pong", func(c *gin.Context) {
		msgs := sessions.Flashes(c)
		if len(msgs) > 0 {
			c.JSON(200, gin.H{
				"message": msgs[0].(string),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "",
			})
		}
	})
	r.Run("0.0.0.0:9090")
}
