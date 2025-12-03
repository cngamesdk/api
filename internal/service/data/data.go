package data

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/log"
	"context"
	"github.com/panjf2000/ants/v2"
	"go.uber.org/zap"
	"time"
)

var (
	launchReportPool *ants.Pool
	gameReportPool   *ants.Pool
)

func init() {
	myLaunchReportPool, launchPoolErr := ants.NewPool(100)
	if launchPoolErr != nil {
		global.Logger.Error("创建协程池异常", zap.Any("err", launchPoolErr))
		return
	}
	launchReportPool = myLaunchReportPool

	myGameReportPool, gameReportPoolErr := ants.NewPool(1000)
	if gameReportPoolErr != nil {
		global.Logger.Error("创建协程池异常", zap.Any("err", gameReportPoolErr))
		return
	}
	gameReportPool = myGameReportPool
}

type DataReportService struct {
}

func (receiver DataReportService) Launch(ctx context.Context, req *api.LaunchReportReq) (resp api.LaunchReportResp, err error) {
	err = launchReportPool.Submit(func() {
		reportModel := log.NewOdsLaunchLogModel()
		reportModel.PlatformId = req.PlatformId
		reportModel.GameId = req.GameId
		reportModel.Action = req.Action
		reportModel.ActionTime = time.Now()
		reportModel.AgentId = req.AgentId
		reportModel.SiteId = req.SiteId
		reportModel.Imei = req.Imei
		reportModel.Idfv = req.Idfv
		reportModel.Oaid = req.Oaid
		reportModel.AndriodId = req.AndriodId
		reportModel.Model = req.Model
		reportModel.Brand = req.Brand
		reportModel.SystemVersion = req.SystemVersion
		reportModel.SdkVersionCode = req.SdkVersionCode
		reportModel.ClientIp = req.ClientIp
		reportModel.Ipv6 = req.Ipv6
		reportModel.Ipv4 = req.Ipv4
		reportModel.UserAgent = req.UserAgent
		if req.MediaSiteId > 0 {
			reportModel.MediaSiteId = req.SiteId
			reportModel.SiteId = req.MediaSiteId
		}
		if saveErr := reportModel.Create(ctx); saveErr != nil {
			global.Logger.Error("保存异常", zap.Any("err", saveErr))
		}
	})
	return
}

func (receiver DataReportService) Game(ctx context.Context, req *api.GameReportReq) (resp api.GameReportResp, err error) {
	err = gameReportPool.Submit(func() {
		reportModel := log.NewOdsGameBehaviorLogModel()
		reportModel.PlatformId = req.PlatformId
		reportModel.UserId = req.UserId
		reportModel.GameId = req.GameId
		reportModel.Action = req.Action
		reportModel.ActionTime = time.Now()
		reportModel.AgentId = req.AgentId
		reportModel.SiteId = req.SiteId
		reportModel.Imei = req.Imei
		reportModel.Idfv = req.Idfv
		reportModel.Oaid = req.Oaid
		reportModel.AndriodId = req.AndriodId
		reportModel.Model = req.Model
		reportModel.Brand = req.Brand
		reportModel.SystemVersion = req.SystemVersion
		reportModel.SdkVersionCode = req.SdkVersionCode
		reportModel.ClientIp = req.ClientIp
		reportModel.Ipv6 = req.Ipv6
		reportModel.Ipv4 = req.Ipv4
		reportModel.ServerId = req.ServerId
		reportModel.ServerName = req.ServerName
		reportModel.RoleId = req.RoleId
		reportModel.RoleName = req.RoleName
		reportModel.RoleLevel = req.RoleLevel
		reportModel.UserAgent = req.UserAgent
		if req.MediaSiteId > 0 {
			reportModel.MediaSiteId = req.SiteId
			reportModel.SiteId = req.MediaSiteId
		}
		if saveErr := reportModel.Create(ctx); saveErr != nil {
			global.Logger.Error("保存异常", zap.Any("err", saveErr))
		}
	})
	return
}
