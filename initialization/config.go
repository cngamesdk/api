package initialization

import (
	"cngamesdk.com/api/global"
	"fmt"
	"github.com/cngamesdk/go-core/cryptor"
	cryptor2 "github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"time"
)

func InitConfigData(myViper *viper.Viper) (err error) {
	if global.Config.Installed > 0 {
		return
	}

	//备份配置文件
	backConfigPath := global.ConfigPath + ".bak." + time.Now().Format("20060102150405")
	myViper.SetConfigFile(backConfigPath)
	if writeErr := myViper.WriteConfig(); writeErr != nil {
		err = errors.Wrap(writeErr, "备份配置文件异常")
		return
	}

	//还原配置文件
	myViper.SetConfigFile(global.ConfigPath)

	backPriKey, backPubKey, backErr := cryptor.GenRsaKey(4096)
	if backErr != nil {
		err = backErr
		return
	}
	myViper.Set("installed", 1)
	myViper.Set("common.background_rsa_private_key", string(backPriKey))
	myViper.Set("common.background_rsa_public_key", string(backPubKey))
	frontPriKey, frontPubKey, frontErr := cryptor.GenRsaKey(4096)
	if frontErr != nil {
		err = frontErr
		return
	}
	myViper.Set("common.front_rsa_private_key", string(frontPriKey))
	myViper.Set("common.front_rsa_public_key", string(frontPubKey))

	myViper.Set("common.sql_md5_crypt_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.aes_crypt_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.token_crypt_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.token_sign_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.common_hash_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.publishing_pay_channel_id", 0)

	myViper.Set("common.game_hash_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	myViper.Set("common.game_hash_key", cryptor2.Md5String(fmt.Sprintf("%s%d%s", random.RandString(10), time.Now().UnixMilli(), random.RandString(10))))

	if writeErr := myViper.WriteConfig(); writeErr != nil {
		err = writeErr
		return
	}

	if readErr := myViper.ReadInConfig(); readErr != nil {
		err = readErr
		return
	}

	if unErr := myViper.Unmarshal(&global.Config); unErr != nil {
		err = unErr
		return
	}

	return
}
