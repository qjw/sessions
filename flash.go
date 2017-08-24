package sessions

import (
	"encoding/gob"
	"github.com/gin-gonic/gin"
	"log"
)

/**
flask用于前后端（不）分离的模型中，不同的model之间相互传递消息

他们通常是使用重定向来关联，不能直接相互通信，所有基于cookie的flash消息是一个不错的主意
 */

var (
	defaultCookieKey       = "_flash"
	defaultStore     Store = nil
)

// 添加flash消息
func AddFlash(c *gin.Context, msg string) {
	if defaultStore == nil{
		log.Fatal("not init yet")
	}
	session, _ := defaultStore.Get(c, defaultCookieKey)
	session.AddFlash(msg)
	session.Save()
}

// 获取所有的flask，并且清空。
func Flashes(c *gin.Context) []interface{} {
	if defaultStore == nil{
		log.Fatal("not init yet")
	}
	session, _ := defaultStore.Get(c, defaultCookieKey)
	flashes := session.Flashes()
	defaultStore.Delete(c, defaultCookieKey)
	return flashes
}

func InitFlash(keyPairs []byte) bool {
	if defaultStore == nil {
		store := NewCookieStore(keyPairs)
		// 不要设置立即失效
		// store.MaxAge(-1)

		defaultStore = store
		gob.Register([]interface{}{})
		return true
	} else {
		log.Fatal("init more than once")
		return false
	}
}
