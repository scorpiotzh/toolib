package toolib

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func JwtString(jwtKey string, expired time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Audience:  "",                                             //用户
		ExpiresAt: time.Now().Unix() + int64(expired/time.Second), //到期时间
		Id:        "",                                             //jwt标识
		IssuedAt:  time.Now().Unix(),                              //发布时间
		Issuer:    "",                                             //发行人
		NotBefore: time.Now().Unix(),                              //在此之前不可用
		Subject:   "",                                             //主题
	})

	return token.SignedString([]byte(jwtKey))
}

func JwtVerify(tokenString, jwtKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("jwt parse:%v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		return nil, fmt.Errorf("verify fail")
	} else {
		return claims, nil
	}
}
