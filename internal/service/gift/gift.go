package gift

import (
	"cngamesdk.com/api/global"
	"cngamesdk.com/api/model/api"
	"cngamesdk.com/api/model/sql/gift"
	"context"
	"fmt"
	gift2 "github.com/cngamesdk/go-core/model/sql/gift"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"time"
)

type GiftService struct {
}

// List 礼包列表
func (receiver GiftService) List(ctx context.Context, req *api.GiftListReq) (resp api.GiftListResp, err error) {
	giftListModel := gift.NewOdsGiftListModel()
	giftCodeClaimModel := gift.NewOdsGiftClaimLogModel()
	tempDb := giftListModel.Db().WithContext(ctx).Table(giftListModel.TableName() + " as gift").
		Joins("LEFT join " + giftCodeClaimModel.TableName() + " code on gift.id = code.gift_id")
	var count int64
	if countErr := tempDb.Count(&count).Error; countErr != nil {
		err = countErr
		global.Logger.Error("获取总数异常", zap.Any("err", countErr))
		return
	}
	var list []api.GiftListDataResp
	if listErr := tempDb.
		Select(fmt.Sprintf("gift.id as id,gift.icon, gift.title, gift.desc, gift.introduce,if(code.code != '','%s',gift.status) as status, ifnull(code.code, '') as code", gift2.GiftCodeStatusClaimed)).
		Limit(req.PageSize).Offset((req.Page - 1) * req.PageSize).
		Find(&list).Error; listErr != nil {
		err = listErr
		global.Logger.Error("获取列表异常", zap.Any("err", listErr))
		return
	}
	resp.Page = req.Page
	resp.PageSize = req.PageSize
	resp.TotalPage = cast.ToInt(math.Ceil(cast.ToFloat64(count) / cast.ToFloat64(req.PageSize)))
	resp.TotalRecord = cast.ToInt(count)
	resp.List = list
	return
}

// Claim 领取
func (receiver GiftService) Claim(ctx context.Context, req *api.GiftClaimReq) (resp api.GiftClaimResp, err error) {
	giftListModel := gift.NewOdsGiftListModel()
	if takeErr := giftListModel.Take(ctx, "*", "id = ?", req.GiftId); takeErr != nil {
		err = takeErr
		global.Logger.Error("获取异常", zap.Any("err", takeErr))
		return
	}
	if !giftListModel.Valid() {
		err = errors.New("礼包已经下架或者领完啦")
		return
	}

	code := ""
	//开启事务
	transactionErr := giftListModel.Db().Transaction(func(tx *gorm.DB) error {
		giftCodeModel := gift.NewOdsGiftCodeListModel()
		if takeErr := tx.WithContext(ctx).Table(giftCodeModel.TableName()).
			Select("id,code").
			Where("gift_id = ? and status = ?", req.GiftId, gift2.GiftCodeStatusNormal).
			Take(giftCodeModel).Error; takeErr != nil {
			global.Logger.Error("获取异常", zap.Any("err", takeErr))
			return takeErr
		}

		updateGiftCodeModel := gift.NewOdsGiftCodeListModel()
		updateGiftCodeModel.Status = gift2.GiftCodeStatusClaimed
		if updateErr := tx.WithContext(ctx).Table(updateGiftCodeModel.TableName()).
			Where("id = ?", giftCodeModel.Id).
			Updates(updateGiftCodeModel).Error; updateErr != nil {
			global.Logger.Error("更新异常", zap.Any("err", updateErr))
			return updateErr
		}
		if updateErr := tx.WithContext(ctx).
			Table(giftListModel.TableName()).
			Where("id = ?", req.GiftId).
			UpdateColumn("available_num", gorm.Expr("available_num - ?", 1)).Error; updateErr != nil {
			global.Logger.Error("更新计数异常", zap.Any("err", updateErr))
			return updateErr
		}
		//写日志
		claimLogModel := gift.NewOdsGiftClaimLogModel()
		claimLogModel.PlatformId = req.PlatformId
		claimLogModel.UserId = req.UserId
		claimLogModel.ActionTime = time.Now()
		claimLogModel.GiftId = req.GiftId
		claimLogModel.Code = giftCodeModel.Code
		if createErr := tx.WithContext(ctx).Table(claimLogModel.TableName()).Create(claimLogModel).Error; createErr != nil {
			global.Logger.Error("保存日志异常", zap.Any("err", createErr))
			return createErr
		}
		code = giftCodeModel.Code
		return nil
	})

	if transactionErr != nil {
		err = transactionErr
		global.Logger.ErrorCtx(ctx, "事务异常", zap.Any("err", transactionErr))
		return
	}
	resp.Code = code
	return
}
