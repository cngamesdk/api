package gift

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/gift"
	"gorm.io/gorm"
)

// OdsGiftListModel 礼包列表
type OdsGiftListModel struct {
	gift.OdsGiftListModel
}

func NewOdsGiftListModel() *OdsGiftListModel {
	model := &OdsGiftListModel{}
	model.OdsGiftListModel.Db = func() *gorm.DB {
		return model.Db()
	}
	return model
}

func (receiver *OdsGiftListModel) Db() *gorm.DB {
	return global.MyDb
}
