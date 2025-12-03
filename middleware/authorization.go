package middleware

import (
	"cngamesdk.com/api/global"
	token2 "cngamesdk.com/api/internal/service/token"
	error2 "cngamesdk.com/api/model/api/error"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/code"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/gin-gonic/gin"
)

// Authorization 授权
func Authorization() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(global.Config.Common.AuthorizationHeadKey)
		if token == "" {
			response2.Out(c, response.NewGlobalResp().SetCode(code.CodeParamEmpty).SetMsg("授权码为空"))
			c.Abort()
			return
		}
		tokenService := &token2.TokenService{}
		verifyData, verifyErr := tokenService.VerifyToken(c, token)
		if verifyErr != nil {
			response2.Out(c, response.NewGlobalResp().SetCode(error2.CodeTokenExpired).SetMsg(verifyErr.Error()))
			c.Abort()
			return
		}
		c.Set(global.Config.Common.CtxTokenDataKey, verifyData)
		c.Next()
	}
}
