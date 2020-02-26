package handler

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"time"

	pb "github.com/vegchic/fullbottle/auth/proto/auth"
)

var AppSecret = config.C().App.Secret

type Claims struct {
	jwt.StandardClaims
	Uid int64
	IP  string
}

type JwtHandler struct{}

func (a *JwtHandler) GenerateJwtToken(ctx context.Context, req *pb.GenerateJwtTokenRequest, resp *pb.GenerateJwtTokenResponse) error {
	var clientIp string
	if ip, ok := metadata.Get(ctx, "ip"); ok {
		clientIp = ip
	}

	expireTime := time.Now().Unix() + req.Expire
	claims := Claims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireTime,
			Id:        fmt.Sprint(req.GetUserId()),
			IssuedAt:  time.Now().Unix(),
			Issuer:    config.AppIss,
		},
		Uid: req.GetUserId(),
		IP:  clientIp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		return errors.New(config.AuthSrvName, "Signing token failed", common.JwtError)
	}

	resp.Token = signed
	return nil
}

func (a *JwtHandler) ParseJwtToken(ctx context.Context, req *pb.ParseJwtTokenRequest, resp *pb.ParseJwtTokenResponse) error {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(req.GetToken(), &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(AppSecret), nil
	})
	if err != nil {
		return errors.New(config.AuthSrvName, "Invalid jwt token", common.JwtError)
	}

	if !token.Valid {
		return errors.New(config.AuthSrvName, "Invalid jwt token", common.JwtError)
	}

	// check ip
	if ip, ok := metadata.Get(ctx, "ip"); ok {
		if ip != claims.IP {
			return errors.New(config.AuthSrvName, "Invalid IP", common.JwtError)
		}
	}

	resp.UserId = claims.Uid
	return nil
}
