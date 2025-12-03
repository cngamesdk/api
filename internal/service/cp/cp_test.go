package cp

import (
	"fmt"
	"github.com/cngamesdk/go-core/model/sql/common"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/netutil"
	"github.com/spf13/cast"
	"testing"
	"time"
)

func TestGame(t *testing.T) {

	gameId := int64(1)
	userId := int64(2)
	platformId := int64(1)
	gameHashKey := "xxx"
	aesKey := "xxx"

	platform := cryptor.Base64StdEncode(string(cryptor.AesEcbEncrypt([]byte(cast.ToString(platformId)), []byte(aesKey))))

	timestamp := time.Now().Unix()
	token := "xxx"

	loginKey := common.GetGameLoginKey(gameId, gameHashKey)
	signStr := fmt.Sprintf("%s%d%d%s%d%s", platform, gameId, userId, token, timestamp, loginKey)
	println(signStr)
	sign := cryptor.Md5String(signStr)
	authReq := map[string]interface{}{
		"platform":  platform,
		"game_id":   gameId,
		"user_id":   userId,
		"token":     token,
		"timestamp": timestamp,
		"sign":      sign,
	}
	println("\n")
	println(netutil.ConvertMapToQueryString(authReq))
	println("\n")
}
