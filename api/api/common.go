package api

import (
	api2 "cngamesdk.com/api/internal/logic/api"
	"cngamesdk.com/api/model/api"
	response2 "cngamesdk.com/api/model/api/response"
	"github.com/cngamesdk/go-core/model/code"
	"github.com/cngamesdk/go-core/model/response"
	"github.com/cngamesdk/go-core/translator"
	"github.com/gin-gonic/gin"
)

// Init 初始化
func Init(ctx *gin.Context) {
	var req api.InitReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).Init(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountAllocate 账号分配
func AccountAllocate(ctx *gin.Context) {
	var req api.AccountAllocateReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountAllocate(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountReg 账号注册
func AccountReg(ctx *gin.Context) {
	var req api.AccountRegReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountReg(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountLogin 账号登录
func AccountLogin(ctx *gin.Context) {
	var req api.AccountLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountLogin(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// SmsSend 短信发送
func SmsSend(ctx *gin.Context) {
	var req api.SmsSendReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).SmsSend(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// SmsLogin 短信验证码登录
func SmsLogin(ctx *gin.Context) {
	var req api.SmsLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).SmsLogin(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// TokenLogin TOKEN登录
func TokenLogin(ctx *gin.Context) {
	var req api.TokenLoginReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).TokenLogin(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// PasswordRetrievePhone 通过手机找回密码
func PasswordRetrievePhone(ctx *gin.Context) {
	var req api.PasswordRetrievePhoneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).PasswordRetrievePhone(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountRetrievePhone 通过手机找回账号
func AccountRetrievePhone(ctx *gin.Context) {
	var req api.AccountRetrievePhoneReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountRetrievePhone(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// Pay 支付
func Pay(ctx *gin.Context) {
	var req api.PayReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).Pay(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()).SetData(resp))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// PaymentStatusInquiryReq 支付状态查询
func PaymentStatusInquiryReq(ctx *gin.Context) {
	var req api.PaymentStatusInquiryReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).PaymentStatusInquiry(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// PasswordModify 密码修改
func PasswordModify(ctx *gin.Context) {
	var req api.PasswordModifyReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).PasswordModify(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// PhoneBind 手机绑定
func PhoneBind(ctx *gin.Context) {
	var req api.PhoneBindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).PhoneBind(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// IdCardBind 身份证绑定
func IdCardBind(ctx *gin.Context) {
	var req api.IdCardBindReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).IdCardBind(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountDelete 用户注销
func AccountDelete(ctx *gin.Context) {
	var req api.AccountDeleteReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountDelete(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// AccountLogout 用户账号登出
func AccountLogout(ctx *gin.Context) {
	var req api.AccountLogoutReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).AccountLogout(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}

// Heartbeat 心跳
func Heartbeat(ctx *gin.Context) {
	var req api.HeartbeatReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(err).Error()))
		return
	}
	req.Format(ctx)
	resp, respErr := (&api2.CommonLogic{}).Heartbeat(ctx, &req)
	if respErr != nil {
		response2.Out(ctx, response.NewGlobalResp().SetMsg(translator.DealErr(respErr).Error()))
		return
	}
	response2.Out(ctx, response.NewGlobalResp().SetCode(code.CodeSuccess).SetData(resp))
	return
}
