package handler

import (
	"FullBottle/common"
	"FullBottle/config"
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/micro/go-micro/v2/util/log"
	"strconv"
	"time"

	pb "FullBottle/auth/proto/auth"
)

var AppSecret = config.GetSingleConfig("app", "secret")
var AppIss = "github.com/vegchic/FullBottle"


type AuthHandler struct{}

func (a *AuthHandler) GenerateJwtToken(ctx context.Context, req *pb.GenerateJwtTokenRequest, resp *pb.GenerateJwtTokenResponse) error {
	expireTime := time.Now().Unix() + req.Expire
	claims := jwt.StandardClaims{
		ExpiresAt: expireTime,
		IssuedAt:  time.Now().Unix(),
		Id:        fmt.Sprintf("%d", req.GetUserId()),
		Issuer:    AppIss,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString([]byte(AppSecret))
	if err != nil {
		resp.Code, resp.Msg = common.JwtError, "Signing token failed"
		log.Fatal(err)

	}
	resp.Token = signed
	resp.Code, resp.Msg = common.Success, "Success"
	return nil
}

func (a *AuthHandler) ParseJwtToken(ctx context.Context, req *pb.ParseJwtTokenRequest, resp *pb.ParseJwtTokenResponse) error {
	claims := jwt.StandardClaims{}
	token, err := jwt.ParseWithClaims(req.GetToken(), &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(AppSecret), nil
	})
	if err != nil {
		resp.Code, resp.Msg = common.JwtError, err.Error()
		return nil
	}

	if !token.Valid {
		resp.Code, resp.Msg = common.JwtError, "Invalid token"
		return nil
	}

	uid, err := strconv.Atoi(claims.Id)
	if err != nil {
		log.Info(err)
		resp.Code, resp.Msg = common.InternalError, "Invalid token"
		return nil
	}

	resp.UserId = int64(uid)
	resp.Code, resp.Msg = common.Success, "Success"
	return nil
}
