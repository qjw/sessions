package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

func main() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       3,
	})
	if err := redisClient.Ping().Err(); err != nil {
		log.Fatal("failed to connect redis")
	}

	store, err := sessions.NewRediStore(redisClient, []byte("abcdefg"))
	if err != nil {
		log.Print(err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		// 设置session。每个session都包含若干key/value对
		session, _ := store.Get(c, "session_test")
		session.Set("key", "value")
		// 保存
		session.Save()
		// 或者 保存所有的session
		// sessions.Save(c)

		c.Redirect(http.StatusFound, "/pong")
	})

	r.GET("/pong", func(c *gin.Context) {
		// 获取session的值
		session, _ := store.Get(c, "session_test")
		value := session.Get("key")
		if value != nil {
			c.JSON(200, gin.H{
				"message": value.(string),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "",
			})
		}
	})

	r.GET("/middle",
		sessions.GinSessionMiddleware(store,"session_test"),
		func(c *gin.Context) {
			// 使用中间件，自动设置session到gin.Context中，避免大量的全局变量传递
			session := c.MustGet("session").(sessions.Session)
			value := session.Get("key")
			if value != nil {
				c.JSON(200, gin.H{
					"message": value.(string),
				})
			} else {
				c.JSON(200, gin.H{
					"message": "",
				})
			}
		})
	r.Run("0.0.0.0:9090")
}
