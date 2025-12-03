package error

import "github.com/pkg/errors"

var (
	ErrorTokenExpired            = errors.New("TOKEN已经失效")
	ErrorUserStatusValid         = errors.New("用户状态异常")
	ErrorUserNameOrPasswordValid = errors.New("用户名或者密码无效")
)

const (
	CodeTokenExpired = 200000
)
