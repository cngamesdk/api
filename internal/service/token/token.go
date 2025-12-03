package token

import (
	"cngamesdk.com/api/global"
	error2 "cngamesdk.com/api/model/api/error"
	"cngamesdk.com/api/model/sql/user"
	"context"
	"errors"
	"fmt"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/golang-jwt/jwt/v5"
	errors2 "github.com/pkg/errors"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"time"
)

type TokenService struct {
	UserId    int64
	Version   int64 //用户版本
	userModel *user.OdsUserInfoLogModel
}

func NewTokenService(UserId int64, version int64) *TokenService {
	return &TokenService{
		UserId:  UserId,
		Version: version,
	}
}

func (receiver *TokenService) GetUserModel() *user.OdsUserInfoLogModel {
	return receiver.userModel
}

func (receiver *TokenService) getUserTokenCacheKey() string {
	return "user-token-" + cast.ToString(receiver.UserId)
}

func (receiver *TokenService) DelCacheToken(ctx context.Context) (err error) {
	cacheKey := receiver.getUserTokenCacheKey()
	_, actionErr := global.MyRedis.Del(ctx, cacheKey).Result()
	if actionErr != nil {
		err = actionErr
		global.Logger.Error("删除缓存异常", zap.Error(actionErr), zap.Any("data", cacheKey))
		return
	}
	return
}

// Generate 生成token
func (receiver *TokenService) Generate(ctx context.Context, extraReq map[string]interface{}) (resp string, err error) {
	if receiver.UserId <= 0 {
		err = errors.New("user_id为空")
		return
	}
	if receiver.Version <= 0 {
		err = errors.New("用户版本为空")
		return
	}
	expireDuration := time.Hour * 24 * 90 // 90天
	mapClaims := jwt.MapClaims{
		"exp":     time.Now().Add(expireDuration).Unix(),
		"user_id": receiver.UserId,
		"version": receiver.Version,
		"sign":    cryptor.Md5String(cast.ToString(receiver.UserId) + global.Config.Common.TokenSignKey),
	}
	if extraReq != nil {
		for itemKey, itemValue := range extraReq {
			mapClaims[itemKey] = itemValue
		}
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	tokenResp, tokenErr := token.SignedString([]byte(global.Config.Common.TokenCryptKey))
	if tokenErr != nil {
		err = tokenErr
		global.Logger.Error("生成TOKEN异常", zap.Any("err", tokenErr))
		return
	}

	cacheKey := receiver.getUserTokenCacheKey()
	_, setErr := global.MyRedis.Set(ctx, cacheKey, tokenResp, expireDuration).Result()
	if setErr != nil {
		err = setErr
		global.Logger.Error("设置TOKEN异常", zap.Any("err", setErr), zap.Any("data", cacheKey))
		return
	}
	resp = tokenResp
	return
}

// VerifyTokenSimple 纯校验jWT有效性
func (receiver *TokenService) VerifyTokenSimple(ctx context.Context, jwtStr string) (resp map[string]interface{}, err error) {
	token, tokenErr := jwt.Parse(jwtStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(global.Config.Common.TokenCryptKey), nil
	}, jwt.WithExpirationRequired())
	if tokenErr != nil {
		err = tokenErr
		return
	}
	// 校验 Claims 对象是否有效，基于 exp（过期时间），nbf（不早于），iat（签发时间）等进行判断（如果有这些声明的话）。
	if !token.Valid {
		err = errors2.Wrap(error2.ErrorTokenExpired, "valid")
		return
	}
	mapClaims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("claims格式错误")
		return
	}
	claimSign, claimExist := mapClaims["sign"]
	if !claimExist {
		err = errors.New("sign字段不存在")
		return
	}
	userId, userIdExists := mapClaims["user_id"]
	if !userIdExists {
		err = errors.New("用户ID不存在")
		return
	}
	receiver.UserId = cast.ToInt64(userId)
	if receiver.UserId <= 0 {
		err = errors.New("用户ID为空")
		return
	}
	mySign := cryptor.Md5String(cast.ToString(receiver.UserId) + global.Config.Common.TokenSignKey)
	if claimSign != mySign {
		err = errors.New("sign加密验证失败")
		return
	}

	resp = mapClaims
	return
}

// VerifyToken TOKEN校验+单点登录校验
func (receiver *TokenService) VerifyToken(ctx context.Context, jwtStr string) (resp map[string]interface{}, err error) {
	verifyTokenSimpleResult, verifyTokenSimpleErr := receiver.VerifyTokenSimple(ctx, jwtStr)
	if verifyTokenSimpleErr != nil {
		err = verifyTokenSimpleErr
		global.Logger.Error("校验TOKEN异常", zap.Any("err", verifyTokenSimpleErr))
		return
	}

	//单点登录校验
	cacheKey := receiver.getUserTokenCacheKey()
	cacheToken, cacheErr := global.MyRedis.Get(ctx, cacheKey).Result()
	if cacheErr != nil {
		err = cacheErr
		global.Logger.Error("获取缓存异常", zap.Any("err", cacheErr), zap.Any("data", cacheKey))
		return
	}
	if jwtStr != cacheToken {
		err = error2.ErrorTokenExpired
		global.Logger.Warn("TOKEN失效", zap.Any("err", err), zap.Any("data", fmt.Sprintf("cache:%s,req:%s", cacheToken, jwtStr)))
		return
	}

	resp = verifyTokenSimpleResult
	return
}

// Verify JWT校验 + 单点登录校验 + 用户校验
func (receiver *TokenService) Verify(ctx context.Context, jwtStr string) (resp map[string]interface{}, err error) {
	verifyResp, verifyErr := receiver.VerifyToken(ctx, jwtStr)
	if verifyErr != nil {
		err = verifyErr
		return
	}
	versionI, ok := verifyResp["version"]
	if !ok {
		err = errors.New("版本不存在")
		return
	}
	version := cast.ToInt64(versionI)
	if version <= 0 {
		err = errors.New("版本格式不正确")
		return
	}
	userModel := user.NewOdsUserInfoLogModel()
	userInfoErr := userModel.Take(ctx, "*", "id = ?", receiver.UserId)
	if userInfoErr != nil {
		err = userInfoErr
		return
	}
	if !userModel.Valid() {
		err = error2.ErrorUserStatusValid
		return
	}
	if userModel.Version > version {
		err = errors2.Wrap(error2.ErrorTokenExpired, "version")
		return
	}
	resp = verifyResp
	receiver.userModel = userModel
	return
}
