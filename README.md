# Session

目前支持三个session后端

1. cookie，session的内容全部序列化到cookie中返回到浏览器，Flash使用此方式
2. file，session的内容存在**本地文件**中，session的id通过cookie返回到浏览器
3. redis，session的内容存在**redis数据库**中，session的id通过cookie返回到浏览器

很少直接使用session

``` go
package main

import (
	"github.com/gin-gonic/gin"
	"github.com/qjw/session"
	"gopkg.in/redis.v5"
	"log"
	"net/http"
)

func GinSessionMiddleware(store sessions.Store, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c, key)
		c.Set("session", session)
	}
}

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
		store.Save(c, session)
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
		GinSessionMiddleware(store,"session_test"),
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
```

# Flask

由于**[gorilla/securecookie](https://github.com/gorilla/securecookie)**需要一个初始密钥进行加密，所以初始化有个密钥的参数

``` go
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
```

输入<http://127.0.0.1:9090/ping> 自动跳转到<http://127.0.0.1:9090/pong>，并且显示ping设置的"hello world"
