package util

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/vegchic/fullbottle/config"
	"time"
)

var appSecret = config.C().App.Secret

type Claims struct {
	jwt.StandardClaims
	Uid int64
	IP  string
}

func GenerateJwtToken(userId int64, expire int64, ip string) (string, error) {
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expire,
			Id:        fmt.Sprint(userId),
			IssuedAt:  time.Now().Unix(),
			Issuer:    config.AppIss,
		},
		Uid: userId,
		IP:  ip,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(appSecret))
	if err != nil {
		return "", err
	}

	return signed, nil
}

func ParseJwtToken(token string, ip string) (claims *Claims, err error) {
	claims = &Claims{}
	t, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	if err != nil {
		return
	}

	if !t.Valid {
		return nil, errors.New("invalid token")
	}

	// check ip
	if ip != claims.IP {
		return nil, errors.New("invalid ip")
	}

	return
}
