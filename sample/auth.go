package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

type User struct {
	Id   int
	Name string
}

func initStore() sessions.Store{
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
	return store
}

func main() {
	store := initStore()
	r := gin.Default()
	r.Use(sessions.GinSessionMiddleware(store, sessions.AUTH_SESSION_NAME))
	r.Use(sessions.GinAuthMiddleware(&sessions.AuthOptions{
		User:&User{},
	}))

	r.GET("/index",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/login",
		func(c *gin.Context) {
			// 是否已经登录
			if sessions.IsAuthenticated(c){
				c.Redirect(http.StatusFound, "/index")
				return
			}
			// 登录授权
			sessions.Login(c,&User{
				Id:1,
				Name:"king",
			})
			c.Redirect(http.StatusFound, "/index")
		})
	r.GET("/logout",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 注销登录
			sessions.Logout(c)
			c.JSON(http.StatusFound, "/logout")
		})
	r.Run("0.0.0.0:9090")
}
