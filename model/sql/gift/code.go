package gift

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/gift"
	"gorm.io/gorm"
)

// OdsGiftCodeListModel 礼包码列表
type OdsGiftCodeListModel struct {
	gift.OdsGiftCodeListModel
}

func NewOdsGiftCodeListModel() *OdsGiftCodeListModel {
	model := &OdsGiftCodeListModel{}
	model.OdsGiftCodeListModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsGiftCodeListModel) Db() *gorm.DB {
	return global.MyDb
}
