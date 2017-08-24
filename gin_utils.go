package sessions

import "github.com/gin-gonic/gin"

func GinSessionMiddleware(store Store, key string) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, _ := store.Get(c, key)
		c.Set("session", session)
	}
}
