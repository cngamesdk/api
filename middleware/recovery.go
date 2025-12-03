package middleware

import (
	"cngamesdk.com/api/global"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/code"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Recovery() func(c *gin.Context) {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil { //用于捕获panic
				global.Logger.ErrorCtx(c, "系统异常", zap.Any("err", err))
				response2.Out(c, response.NewGlobalResp().SetCode(code.CodeSystemErr).SetMsg("系统异常"))
				c.Abort()
			}
		}()
		c.Next() // 调用下一个处理
	}
}
