package jwt

import (
	"fmt"
	"github.com/astaxie/beego"
	jwtv3 "gopkg.in/dgrijalva/jwt-go.v3"
)

//MyCustomClaims 自定义信息
type MyWxAppClaims struct {
	UserId     string `json:"user_id"`
	SessionKey string `json:"session_key"`
	OpenId     string `json:"open_id"`
	jwtv3.StandardClaims
}

//WxClaimsToken 微信app token加密
func WxClaimsToken(cl *MyWxAppClaims, secret string) (string, error) {
	token := jwtv3.NewWithClaims(jwtv3.SigningMethodHS256, cl)
	tokenString, err := token.SignedString([]byte(secret))
	return tokenString, err
}

//WxClaimsParse 微信小程序信息解析
func WxClaimsParse(token, secret string) (jwtv3.MapClaims, error) {
	var nilClaims jwtv3.MapClaims
	t, err := jwtv3.Parse(token, func(*jwtv3.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		beego.Error(fmt.Sprintf("解密出现错误%v，token：%s， secret：%s", err, token, secret))
		return nilClaims, err
	}

	claims, ok := t.Claims.(jwtv3.MapClaims)
	if !ok {
		beego.Error(fmt.Sprintf("用户认证失败：%v, claims:%#v", err, t.Claims))
		return nilClaims, err
	}
	return claims, err
}
