package initialization

import (
	"cngamesdk.com/api/global"
	"github.com/cngamesdk/go-core/model/sql/advertising"
	common2 "github.com/cngamesdk/go-core/model/sql/common"
	gift2 "github.com/cngamesdk/go-core/model/sql/gift"
	log2 "github.com/cngamesdk/go-core/model/sql/log"
	publishing2 "github.com/cngamesdk/go-core/model/sql/publishing"
	user2 "github.com/cngamesdk/go-core/model/sql/user"
)

// Migrate 迁移数据
func Migrate() (err error) {
	if migrateErr := global.MyDb.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci ROW_FORMAT=DYNAMIC").
		AutoMigrate(
			&common2.DimPlatformModel{},
			&common2.DimCompanyModel{},
			&common2.DimGameModel{},
			&common2.DimMainGameModel{},
			&common2.DimRootGameModel{},
			&common2.DimPayChannelModel{},
			&common2.DimProductCommonConfigurationModel{},
			&common2.DimGameAppVersionConfiguration{},
			&common2.DimAgentModel{},
			&common2.DimSiteModel{},
			&user2.OdsUserSmsSendLogModel{},
			&user2.OdsUserInfoLogModel{},
			&user2.OdsUserOperationLogModel{},
			&gift2.OdsGiftCodeListModel{},
			&gift2.OdsGiftListModel{},
			&gift2.OdsGiftClaimLogModel{},
			&log2.OdsGameBehaviorLogModel{},
			&log2.DwdGameRegLogModel{},
			&log2.OdsLaunchLogModel{},
			&log2.OdsLoginLogModel{},
			&log2.OdsPayLogModel{},
			&log2.OdsRegLogModel{},
			&log2.DwdRootGameBackRegLogModel{},
			&log2.DwdRootGameRegLogModel{},
			&publishing2.DimPublishingChannelConfigModel{},
			&publishing2.DimPublishingChannelGameConfigModel{},
			&publishing2.OdsPublishingUserBindLogModel{},
			&advertising.OdsMediaCallbackLogModel{},
			&advertising.OdsMediaAdClickLogModel{},
		); migrateErr != nil {
		err = migrateErr
	}
	return
}
