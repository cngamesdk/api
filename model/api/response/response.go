package response

import (
	"cngamesdk.com/api/global"
	"fmt"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func Out(ctx *gin.Context, data response.GlobalResp) {
	requestId := data.RequestId
	if requestId == "" {
		ctxRequestId := ctx.GetString(global.Config.Common.CtxRequestIdKey)
		requestId = ctxRequestId
		if ctxRequestId == "" {
			tmpUuid, tmpErr := random.UUIdV4()
			if tmpErr != nil {
				tmpUuid = cryptor.Md5String(fmt.Sprintf("%d%s", time.Now().UnixMilli(), random.RandString(5)))
			}
			requestId = tmpUuid
		}
		data.RequestId = requestId
	}
	ctx.JSON(http.StatusOK, data)
}
