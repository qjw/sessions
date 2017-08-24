package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

type User2 struct {
	Id   int
	Name string
}

func initStore2() sessions.Store {
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
	store := initStore2()
	r := gin.Default()
	r.Use(sessions.GinSessionMiddleware(store, sessions.AUTH_SESSION_NAME))
	r.Use(sessions.GinAuthMiddleware(&sessions.AuthOptions{
		User: &User2{},
	}))
	sessions.InitPermission(&sessions.PermissionOptions{
		UserPermissionGetter: func(user interface{}) (map[int]bool, error) {
			ruser := user.(*User2)
			if ruser.Name == "p1" {
				return map[int]bool{
					1: true,
				}, nil
			} else if ruser.Name == "p2" {
				return map[int]bool{
					1: true,
					2: true,
				}, nil
			} else {
				return map[int]bool{}, nil
			}
		},
		AllPermisionsGetter: func() (map[string]int, error) {
			return map[string]int{
				"perm1": 1,
				"perm2": 2,
				"perm3": 3,
			}, nil
		},
	})

	r.GET("/index",
		sessions.LoginRequired(),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm1",
		sessions.PermissionRequired("perm1"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm2",
		sessions.PermissionRequired("perm2"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/perm3",
		sessions.PermissionRequired("perm3"),
		func(c *gin.Context) {
			// 获取登录用户
			user := sessions.LoggedUser(c).(*User2)
			c.JSON(http.StatusOK, gin.H{
				"message": user.Name,
			})
		})
	r.GET("/login",
		func(c *gin.Context) {
			// 是否已经登录
			if sessions.IsAuthenticated(c) {
				c.Redirect(http.StatusFound, "/index")
				return
			}

			// 登录授权
			sessions.Login(c, &User2{
				Id:   1,
				Name: c.DefaultQuery("name", "p1"),
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
