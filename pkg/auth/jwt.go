package auth

import (
	"errors"
	"time"
)

import "github.com/golang-jwt/jwt"

var secret = []byte("4df17f1a-9157-4cd4-ac18-47c42ba374fd")

type Info struct {
	UID string
	//OrgID          string
	IsRefreshToken bool
}

type JWTClaims struct {
	Info Info
	jwt.StandardClaims
}

// 上线根据需求改
const (
	AccessTokenExpireIn  = time.Hour * 24
	RefreshTokenExpireIn = time.Hour * 24 * 30
)

// GenToken 生成JWT
func GenToken(info Info, expire ...time.Duration) (token string, err error) {
	if len(expire) == 0 {
		expire = append(expire, AccessTokenExpireIn)
	}
	c := JWTClaims{
		Info: info,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(expire[0]).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "05sec",
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	return t.SignedString(secret)
}

// ParseToken 解析JWT
func ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid { // 校验token
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
