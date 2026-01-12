package pay_switch

import (
	"cngamesdk.com/api/global"
	pay2 "cngamesdk.com/api/internal/service/pay"
	"cngamesdk.com/api/internal/service/pay/channel"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/cache/game"
	"cngamesdk.com/api/model/cache/pay"
	"cngamesdk.com/api/model/sql/log"
	"context"
	"fmt"
	"github.com/cngamesdk/go-core/model/sql"
	"github.com/cngamesdk/go-core/model/sql/common"
	"github.com/duke-git/lancet/v2/mathutil"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"slices"
)

const (
	OperatorIn       = "in"
	OperatorNotIn    = "not-in"
	OperatorEqual    = "equal"
	OperatorNotEqual = "not-equal"
)

var defaultRules = map[string]func(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (bool, error){
	"platform_id":  rulePlatformId,
	"root_game_id": ruleRootGameId,
	"main_game_id": ruleMainGameId,
	"game_id":      ruleGameId,
	"agent_id":     ruleAgentId,
	"site_id":      ruleSiteId,
	"user_id":      ruleUserId,
	"oaid":         ruleOaid,
	"imei":         ruleImei,
	"idfv":         ruleIdfv,
}

func handleRule(value interface{}, rule common.DimPayChannelSwitchRule) (resp bool) {
	switch rule.Operator {
	case OperatorIn:
		resp = slices.Contains(rule.Value, value)
		return
	case OperatorNotIn:
		resp = !slices.Contains(rule.Value, value)
		return
	case OperatorEqual:
		if len(rule.Value) <= 0 {
			resp = false
			return
		}
		resp = value == rule.Value[0]
		return
	case OperatorNotEqual:
		if len(rule.Value) <= 0 {
			resp = false
			return
		}
		resp = value != rule.Value[0]
		return
	default:
		return
	}
}

// rulePlatformId 平台规则
func rulePlatformId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.PlatformId, rule)
	return
}

// ruleRootGameId 根游戏ID规则
func ruleRootGameId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	gameModel := game.NewDimGameModel()
	gameDetailErr := gameModel.DetailInfoByGameId(ctx, req.GameId)
	if gameDetailErr != nil {
		err = gameDetailErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Error(gameDetailErr))
		return
	}
	resp = handleRule(gameModel.RootGameId, rule)
	return
}

// ruleMainGameId 主游戏ID规则
func ruleMainGameId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	gameModel := game.NewDimGameModel()
	gameDetailErr := gameModel.DetailInfoByGameId(ctx, req.GameId)
	if gameDetailErr != nil {
		err = gameDetailErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Error(gameDetailErr))
		return
	}
	resp = handleRule(gameModel.MainGameId, rule)
	return
}

// ruleGameId 子游戏ID规则
func ruleGameId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.GameId, rule)
	return
}

// ruleAgentId 渠道ID规则
func ruleAgentId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	model := log.NewDwdGameRegLogModel()
	if takeErr := model.Take(ctx, "*", "platform_id = ? and game_id = ? and user_id = ?", req.PlatformId, req.GameId, req.UserId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Error(takeErr))
		return
	}
	resp = handleRule(model.AgentId, rule)
	return
}

// ruleSiteId 广告位ID规则
func ruleSiteId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	model := log.NewDwdGameRegLogModel()
	if takeErr := model.Take(ctx, "*", "platform_id = ? and game_id = ? and user_id = ?", req.PlatformId, req.GameId, req.UserId); takeErr != nil {
		err = takeErr
		global.Logger.ErrorCtx(ctx, "获取异常", zap.Error(takeErr))
		return
	}
	resp = handleRule(model.SiteId, rule)
	return
}

// ruleUserId 用户ID规则
func ruleUserId(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.UserId, rule)
	return
}

// ruleOaid 安卓OAID规则
func ruleOaid(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.Oaid, rule)
	return
}

// ruleImei 安卓IMEI/iOS Idfa规则
func ruleImei(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.Imei, rule)
	return
}

// ruleIdfv iOS Idfv规则
func ruleIdfv(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (resp bool, err error) {
	resp = handleRule(req.Idfv, rule)
	return
}

type PayChannelSwitchService struct {
}

func (receiver *PayChannelSwitchService) getRules(ctx context.Context) (
	resp map[string]func(ctx context.Context, req *api.PayReq, rule common.DimPayChannelSwitchRule) (bool, error), err error) {
	resp = defaultRules
	return
}

// getPayChannelByWeight 通过权重获取支付渠道
func (receiver *PayChannelSwitchService) getPayChannelByWeight(ctx context.Context, req pay.DimPayChannelSwitchModel) (resp int64, err error) {
	type payChannelItem struct {
		existsNum                            int    // 已经切的数量
		cacheKey                             string // 缓存键值
		common.DimPayChannelSwitchPayChannel        // 配置
	}
	totalWeight := 0
	totalExistsNum := 0
	cacheKeyFormat := "pay-channel-switch-%d-%d"
	var payChannelsContainer []payChannelItem
	cacheClient := global.MyRedis
	for _, item := range req.PayChannels {
		existsNum := 0
		totalWeight += item.Weight
		tempCacheKey := fmt.Sprintf(cacheKeyFormat, req.Id, item.PayChannelId)
		existsResult, existsErr := cacheClient.Exists(ctx, tempCacheKey).Result()
		if existsErr != nil {
			err = existsErr
			global.Logger.ErrorCtx(ctx, "缓存Exists异常", zap.Error(existsErr))
			return
		}
		if existsResult > 0 {
			getResult, getErr := cacheClient.Get(ctx, tempCacheKey).Result()
			if getErr != nil {
				err = getErr
				global.Logger.ErrorCtx(ctx, "获取缓存异常", zap.Error(getErr))
				return
			}
			existsNum = cast.ToInt(getResult)
		}
		totalExistsNum += existsNum
		payChannelsContainer = append(payChannelsContainer, payChannelItem{cacheKey: tempCacheKey, existsNum: existsNum, DimPayChannelSwitchPayChannel: item})
	}
	for _, item := range payChannelsContainer {
		configRate := mathutil.Percent(cast.ToFloat64(item.Weight), cast.ToFloat64(totalWeight), 5)
		existsRate := mathutil.Percent(cast.ToFloat64(item.existsNum), cast.ToFloat64(totalExistsNum), 5)
		if existsRate <= configRate { //必须等于，防止单个配置

			//设置缓存，自增
			if _, incrErr := cacheClient.Incr(ctx, item.cacheKey).Result(); incrErr != nil {
				err = incrErr
				global.Logger.ErrorCtx(ctx, "设置自增数量异常", zap.Error(incrErr), zap.Any("data", item.cacheKey))
				return
			}
			resp = item.PayChannelId

			global.Logger.InfoCtx(ctx, "比例命中", zap.Any("item", item),
				zap.Any("data", []interface{}{configRate, totalWeight, existsRate, totalExistsNum}))

			return
		}
	}
	return
}

func (receiver *PayChannelSwitchService) Handle(ctx context.Context, req *api.PayReq) (resp channel.PayChannelInterface, err error) {
	if req == nil {
		err = errors.New("支付参数不能为空")
		return
	}
	//获取对应支付方式的所有配置
	cacheModel := &pay.DimPayChannelSwitchModel{}
	allConfigRules, getErr := cacheModel.FindAllRules(ctx, "*", "platform_id = ? and pay_type = ? and status = ?",
		req.PlatformId, req.PayType, sql.StatusNormal)
	if getErr != nil {
		err = getErr
		global.Logger.ErrorCtx(ctx, "获取配置异常", zap.Error(getErr))
		return
	}

	payChannelId := int64(0)

	if len(allConfigRules) > 0 {
		configRules, configRulesErr := receiver.getRules(ctx)
		if configRulesErr != nil {
			err = configRulesErr
			global.Logger.ErrorCtx(ctx, "获取配置规则异常", zap.Error(configRulesErr))
			return
		}
		for _, item := range allConfigRules {
			for _, ruleItem := range item.Rules {
				execRule, configRuleOK := configRules[ruleItem.Key]
				if configRuleOK {
					execResult, execErr := execRule(ctx, req, ruleItem)
					if execErr != nil {
						err = execErr
						global.Logger.ErrorCtx(ctx, "执行规则异常", zap.Error(execErr), zap.Any("data", ruleItem))
						return
					}
					//没有命中规则,则搜索下个规则
					if !execResult {
						break
					}
				}
			}
			//命中规则
			rulePayChannelId, getConfigErr := receiver.getPayChannelByWeight(ctx, item)
			if getConfigErr != nil {
				err = getConfigErr
				global.Logger.ErrorCtx(ctx, "执行规则配置支付渠道异常", zap.Error(getConfigErr), zap.Any("data", item))
				return
			}
			payChannelId = rulePayChannelId
			global.Logger.InfoCtx(ctx, "命中规则", zap.Any("item", item), zap.Any("data", rulePayChannelId))
			break
		}

	}
	//获取默认按主体支付渠道
	if payChannelId <= 0 {
		//读取主体对应的支付渠道
		gameInfoModel := game.NewDimGameModel()
		if getGameInfoErr := gameInfoModel.Take(ctx, "*", "id = ?", req.GameId); getGameInfoErr != nil {
			err = getGameInfoErr
			global.Logger.ErrorCtx(ctx, "获取游戏异常", zap.Error(getGameInfoErr))
			return
		}
		gameCompanyId := gameInfoModel.CompanyId
		payChannelModel := pay.NewDimPayChannelModel()
		payChannelErr := payChannelModel.GetRateLessChannel(ctx, "*", "platform_id = ? and company_id = ? and pay_type = ? and status = ?",
			req.PlatformId, gameCompanyId, req.PayType, sql.StatusNormal)
		if payChannelErr != nil {
			err = errors.Wrap(payChannelErr, "未配置支付渠道")
			global.Logger.ErrorCtx(ctx, "未配置支付渠道", zap.Error(payChannelErr))
			return
		}
		payChannelId = payChannelModel.Id
		global.Logger.InfoCtx(ctx, "命中默认主体规则", zap.Any("data", payChannelModel))
	}
	//获取支付渠道网关
	payChannelGateway := pay2.GetPayChannel(payChannelId)
	if payChannelGateway == nil {
		err = errors.New("未找到支付网关")
		global.Logger.ErrorCtx(ctx, "未找到支付网关", zap.Any("data", payChannelId))
		return
	}
	resp = payChannelGateway
	return
}
