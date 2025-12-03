package middleware

import (
	"bytes"
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api"
	response2 "cngamesdk.com/api/model/api/response"
	"encoding/base64"
	"encoding/json"
	"fmt"
	cryptor3 "github.com/cngamesdk/go-core/cryptor"
	code2 "github.com/cngamesdk/go-core/model/code"
	error3 "github.com/cngamesdk/go-core/model/error"
	response3 "github.com/cngamesdk/go-core/model/response"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"time"
)

// AppCrypt APP加解密
func AppCrypt() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		var readErr error
		body, readErr = io.ReadAll(c.Request.Body)
		var bizErr error
		if readErr != nil {
			bizErr = readErr
			global.Logger.Error("read body from request error:", zap.Error(readErr))
		} else {
			var appReq api.AppReq
			jsonErr := json.Unmarshal(body, &appReq)
			if jsonErr != nil {
				bizErr = jsonErr
				global.Logger.Error("JSON解码失败", zap.Error(jsonErr), zap.ByteString("body", body))
			} else {
				//RSA解密获取AES密码
				if appReq.Key == "" {
					bizErr = errors.Wrap(error3.ErrorParamEmpty, "key")
					global.Logger.Error("KEY为空", zap.Error(bizErr))
				} else {
					aesKey, rsaErr := cryptor3.RsaDecryptBase64([]byte(global.Config.Common.FrontRsaPrivateKey), appReq.Key)
					if rsaErr != nil {
						bizErr = rsaErr
						global.Logger.Error("RSA解密失败", zap.Error(rsaErr), zap.ByteString("body", body))
					} else {
						body = cryptor.AesEcbDecrypt([]byte(cryptor.Base64StdDecode(appReq.Data)), aesKey)
						global.Logger.Info("解密后参数", zap.ByteString("body", body))
					}
				}
			}
		}
		if bizErr != nil {
			response2.Out(c, response3.NewGlobalResp().SetCode(code2.CodeDeCryptErr).SetMsg(error3.ErrorDecrypt.Error()))
			c.Abort()
			return
		}

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// 创建缓冲区来捕获所有写入
		bodyBuffer := &bytes.Buffer{}

		// 临时替换 Writer
		originalWriter := c.Writer
		tempWriter := &bufferedWriter{
			ResponseWriter: originalWriter,
			body:           bodyBuffer,
			skipWrite:      true, // 重要：阻止写入原始响应
		}
		c.Writer = tempWriter

		c.Next()

		// 恢复原始 Writer
		c.Writer = originalWriter

		var appResp api.AppResp
		respData := bodyBuffer.Bytes()
		global.Logger.InfoCtx(c, "明文数据", zap.ByteString("body", respData))
		appResp.RequestId = global.Logger.GetRequestId(c)
		aesKey := cryptor.Md5String(fmt.Sprintf("%d%s", time.Now().UnixMilli(), random.RandString(5)))
		rsaAesKey, encryptErr := cryptor3.RsaEncryptBase64([]byte(global.Config.Common.BackgroundRsaPublicKey), []byte(aesKey))
		if encryptErr != nil {
			global.Logger.ErrorCtx(c, "RSA加密失败", zap.Error(encryptErr), zap.ByteString("body", respData))
			response2.Out(c, response3.NewGlobalResp().SetCode(code2.CodeEnCryptErr).SetMsg(error3.ErrorEncrypt.Error()))
			c.Abort()
		}
		appResp.Key = rsaAesKey
		appResp.Data = base64.StdEncoding.EncodeToString(cryptor.AesEcbEncrypt(respData, []byte(aesKey)))

		jsonResp, jsonRespErr := json.Marshal(appResp)
		if jsonRespErr != nil {
			global.Logger.ErrorCtx(c, "JSON_Marshal失败", zap.Error(jsonRespErr))
			response2.Out(c, response3.NewGlobalResp().SetCode(code2.CodeEnCryptErr).SetMsg(error3.ErrorEncrypt.Error()))
			c.Abort()
		}

		// 清空任何可能已设置的 Header
		c.Writer.Header().Del("Content-Length")
		_, writeErr := c.Writer.Write(jsonResp)
		if writeErr != nil {
			global.Logger.ErrorCtx(c, "写入失败", zap.Error(writeErr), zap.ByteString("body", jsonResp))
			response2.Out(c, response3.NewGlobalResp().SetCode(code2.CodeEnCryptErr).SetMsg(error3.ErrorEncrypt.Error()))
			c.Abort()
		}
	}
}

type bufferedWriter struct {
	gin.ResponseWriter
	body      *bytes.Buffer
	skipWrite bool
}

func (w *bufferedWriter) Write(b []byte) (int, error) {
	if w.skipWrite {
		// 只写入缓冲区，不写入原始响应
		return w.body.Write(b)
	}
	// 正常写入
	return w.ResponseWriter.Write(b)
}
