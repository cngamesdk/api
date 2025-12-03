package middleware

import (
	"cngamesdk.com/api/global"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
)

// Trace 链路追踪
func Trace() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestId := c.Query(global.Config.Common.CtxRequestIdKey)
		if requestId == "" {
			requestId = random.RandString(32)
		}
		c.Set(global.Config.Common.CtxRequestIdKey, requestId)
		c.Next()
	}
}
