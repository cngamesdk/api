package initialization

import (
	"cngamesdk.com/api/api/api"
	"cngamesdk.com/api/api/cp"
	"cngamesdk.com/api/api/pay"
	"cngamesdk.com/api/middleware"
	"github.com/gin-gonic/gin"
)

// RouteInit 路由初始化
func RouteInit(engine *gin.Engine) {
	engine.Use(middleware.Record2File())
	initCommonRoute(engine)
	initAppRoute(engine)
}

// initCommonRoute 初始化通用路由
func initCommonRoute(engine *gin.Engine) {
	//支付渠道回调
	payChannelGroup := engine.Group("pay-channel")
	{
		payChannelGroup.POST("callback/:pay_channel_id", pay.PayChannelCallback)
		payChannelGroup.GET("callback/:pay_channel_id", pay.PayChannelCallback)
	}

	//发行渠道分组
	publishingGroup := engine.Group("publishing")
	{
		payGroup := publishingGroup.Group("pay")
		{
			payGroup.GET("callback/:channel_id", pay.PublishingChannelPayCallback)
		}
	}

	//CP分组
	cpGroup := engine.Group("cp")
	{
		cpGroup.GET("auth", cp.Auth)
	}
}

// initAppRoute 初始化APP路由
func initAppRoute(engine *gin.Engine) {

	//APP路由
	appRoutes := engine.Group("app")
	//appRoutes.Use(middleware.AppCrypt(), middleware.CheckAppSign())
	appRoutes.Use(middleware.ClientMobileGame())
	{
		appRoutes.POST("init", api.Init)
		userGroups := appRoutes.Group("user")
		{
			userGroups.POST("account-allocate", api.AccountAllocate)
			userGroups.POST("account-reg", api.AccountReg)
			userGroups.POST("account-login", api.AccountLogin)
			userGroups.POST("phone-sms-send", api.SmsSend)
			userGroups.POST("phone-login", api.SmsLogin)
			userGroups.POST("token-login", api.TokenLogin)
			userGroups.POST("password-retrieve-phone", api.PasswordRetrievePhone)
			userGroups.POST("account-retrieve-phone", api.AccountRetrievePhone)

			//需要授权分组
			userAuthorizationGroups := userGroups.Group("")
			userAuthorizationGroups.Use(middleware.Authorization())
			{
				userAuthorizationGroups.POST("pay", api.Pay)
				userAuthorizationGroups.POST("pay-status-inquiry", api.PaymentStatusInquiryReq)
				userAuthorizationGroups.POST("password-modify", api.PasswordModify)
				userAuthorizationGroups.POST("phone-bind", api.PhoneBind)
				userAuthorizationGroups.POST("idcard-bind", api.IdCardBind)
				userAuthorizationGroups.POST("account-delete", api.AccountDelete)
				userAuthorizationGroups.POST("account-logout", api.AccountLogout)
				userAuthorizationGroups.POST("heartbeat", api.Heartbeat)
			}
		}
		giftGroups := appRoutes.Group("gift")
		giftGroups.Use(middleware.Authorization())
		{
			giftGroups.POST("list", api.GiftList)
			giftGroups.POST("claim", api.GiftClaim)
		}
		//游戏发行
		publishingGroups := appRoutes.Group("publishing")
		{
			publishingGroups.POST("login", api.PublishingLogin)
			publishingGroups.Use(middleware.Authorization()).POST("pay", api.PublishingPay)
		}
		//数据
		dataGroups := appRoutes.Group("data")
		{
			dataGroups.POST("launch-report", api.LaunchDataReport)
			dataGroups.Use(middleware.Authorization()).POST("game-log-report", api.GameLogDataReport)
		}

		//媒体分组
		mediaGroups := appRoutes.Group("media")
		{
			//广告分组
			advertisingGroups := mediaGroups.Group("advertising")
			{
				advertisingGroups.POST("click", api.ReportReg) // 广告点击
			}

			//上报分组
			reportGroups := mediaGroups.Group("report")
			reportGroups.Use(middleware.Authorization())
			{
				reportGroups.POST("reg", api.ReportReg)
				reportGroups.POST("login", api.ReportLogin)
				reportGroups.POST("pay", api.ReportPay)
				reportGroups.POST("callback", api.ReportCallback)
			}
		}
	}
}
