package myjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	logs "github.com/sirupsen/logrus"
)

const (
	MySigningKey = "liH7uomvOtELRWQsga9nJDVU34Zc25ye" // 私钥，（随机串 加盐）
)

// 生成 jwt token
func CreateToken() (tokenStr string, err error) {

	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(20 * time.Second)), // 过期时间，秒
		NotBefore: jwt.NewNumericDate(time.Now()),                       // 生效时间，秒
		Issuer:    "register-server",                                    // 签发人
		IssuedAt:  jwt.NewNumericDate(time.Now()),                       // 签发时间
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err = token.SignedString([]byte(MySigningKey))
	if err != nil {
		logs.Error(err)
		return
	}
	return tokenStr, err
}

// 解析 token
func ParseToken(tokenStr string) (err error) {
	playload := jwt.RegisteredClaims{} // jwt standard playload
	token, err := jwt.ParseWithClaims(
		tokenStr,
		&playload,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(MySigningKey), nil
		},
	)
	if err != nil {
		logs.Error(err)
		return
	}
	logs.Infof("token struct is: %#v", token)
	logs.Infof("playload struct is: %#v", playload)
	return
}
