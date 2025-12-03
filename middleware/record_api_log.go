package middleware

import (
	"bytes"
	"cngamesdk.com/api/global"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var apiRespPool sync.Pool
var recordLogPool *ants.Pool

func init() {
	apiRespPool.New = func() interface{} {
		return make([]byte, 1024)
	}
	myRecordLogPool, recordLogPoolErr := ants.NewPool(100)
	if recordLogPoolErr != nil {
		global.Logger.Error("创建协程池异常", zap.Any("err", recordLogPoolErr))
		return
	}
	recordLogPool = myRecordLogPool
}

func Record2File() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Method != http.MethodGet {
			var err error
			body, err = io.ReadAll(c.Request.Body)
			if err != nil {
				global.Logger.Error("read body from request error:", zap.Error(err))
			} else {
				c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			}
		} else {
			query := c.Request.URL.RawQuery
			query, _ = url.QueryUnescape(query)
			split := strings.Split(query, "&")
			m := make(map[string]string)
			for _, v := range split {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					m[kv[0]] = kv[1]
				}
			}
			body, _ = json.Marshal(&m)
		}
		record := fileRecordData{
			Ip:     c.ClientIP(),
			Method: c.Request.Method,
			Path:   c.Request.URL.Path,
			Agent:  c.Request.UserAgent(),
			Body:   string(body),
		}

		// 上传文件时候 中间件日志进行裁断操作
		if strings.Contains(c.GetHeader("Content-Type"), "multipart/form-data") {
			if len(record.Body) > 1024 {
				// 截断
				newBody := apiRespPool.Get().([]byte)
				copy(newBody, record.Body)
				record.Body = string(newBody)
				defer apiRespPool.Put(newBody[:0])
			}
		}

		writer := responseBodyWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		now := time.Now()

		c.Next()

		latency := time.Since(now)
		record.ErrorMessage = c.Errors.ByType(gin.ErrorTypePrivate).String()
		record.Status = c.Writer.Status()
		record.Latency = latency.Seconds()
		record.Resp = writer.body.String()

		if strings.Contains(c.Writer.Header().Get("Pragma"), "public") ||
			strings.Contains(c.Writer.Header().Get("Expires"), "0") ||
			strings.Contains(c.Writer.Header().Get("Cache-Control"), "must-revalidate, post-check=0, pre-check=0") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/force-download") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/octet-stream") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/vnd.ms-excel") ||
			strings.Contains(c.Writer.Header().Get("Content-Type"), "application/download") ||
			strings.Contains(c.Writer.Header().Get("Content-Disposition"), "attachment") ||
			strings.Contains(c.Writer.Header().Get("Content-Transfer-Encoding"), "binary") {

			if len(record.Resp) > 1024 {
				// 截断
				newBody := apiRespPool.Get().([]byte)
				copy(newBody, record.Resp)
				record.Body = string(newBody)
				defer apiRespPool.Put(newBody[:0])
			}
		}

		//采用协程池
		if submitErr := recordLogPool.Submit(func() {
			//record.Body = log2.ReplaceStrSensitiveData(record.Body, []byte(global.Config.Common.AesCryptKey))
			//record.Resp = log2.ReplaceStrSensitiveData(record.Resp, []byte(global.Config.Common.AesCryptKey))
			global.Logger.InfoCtx(c, "api_record", zap.Any("log", record))
		}); submitErr != nil {
			global.Logger.Error("协程池提交异常", zap.Error(submitErr))
		}
	}
}

type fileRecordData struct {
	Ip           string  `json:"ip" form:"ip" gorm:"column:ip;comment:请求ip"`                                   // 请求ip
	Method       string  `json:"method" form:"method" gorm:"column:method;comment:请求方法"`                       // 请求方法
	Path         string  `json:"path" form:"path" gorm:"column:path;comment:请求路径"`                             // 请求路径
	Status       int     `json:"status" form:"status" gorm:"column:status;comment:请求状态"`                       // 请求状态
	Latency      float64 `json:"latency" form:"latency" gorm:"column:latency;comment:延迟" swaggertype:"string"` // 延迟
	Agent        string  `json:"agent" form:"agent" gorm:"column:agent;comment:代理"`                            // 代理
	ErrorMessage string  `json:"error_message" form:"error_message" gorm:"column:error_message;comment:错误信息"`  // 错误信息
	Body         string  `json:"body" form:"body" gorm:"type:text;column:body;comment:请求Body"`                 // 请求Body
	Resp         string  `json:"resp" form:"resp" gorm:"type:text;column:resp;comment:响应Body"`                 // 响应Body
}

type responseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
