package middleware

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/common"
	"github.com/gin-gonic/gin"
)

// CheckAppSign 验证APP签名中间件
func ClientMobileGame() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set(global.CtxKeyClient, common.GameTypeMobileGame)
		c.Next()
	}
}
