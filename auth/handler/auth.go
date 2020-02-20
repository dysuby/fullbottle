package handler

import (
	"FullBottle/common"
	"FullBottle/common/log"
	"FullBottle/config"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/v2/errors"
	"strconv"
	"time"

	pb "FullBottle/auth/proto/auth"
)

var AppSecret = config.GetConfig().App.Secret

type AuthHandler struct{}

func (a *AuthHandler) GenerateJwtToken(ctx context.Context, req *pb.GenerateJwtTokenRequest, resp *pb.GenerateJwtTokenResponse) error {
	expireTime := time.Now().Unix() + req.Expire
	claims := jwt.StandardClaims{
		ExpiresAt: expireTime,
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", req.GetUserId()),
		Issuer:    config.AppIss,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		return errors.New(config.AuthSrvName, "Signing token failed", common.JwtError)
	}

	resp.Token = signed
	return nil
}

func (a *AuthHandler) ParseJwtToken(ctx context.Context, req *pb.ParseJwtTokenRequest, resp *pb.ParseJwtTokenResponse) error {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(req.GetToken(), &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(AppSecret), nil
	})
	if err != nil {
		return errors.New(config.AuthSrvName, "Parsing claims failed", common.JwtError)
	}

	if !token.Valid {
		return errors.New(config.AuthSrvName, "Invalid jwt token", common.JwtError)
	}

	uid, err := strconv.Atoi(claims.Id)
	if err != nil {
		log.Errorf(err, "Claims format error: %v", token)
		return errors.New(config.AuthSrvName, "Claims format error", common.InternalError)
	}

	resp.UserId = int64(uid)
	return nil
}
