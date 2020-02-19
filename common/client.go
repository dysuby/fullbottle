package common

import (
	"github.com/micro/go-micro/v2/client"

	PbUser "FullBottle/user/proto/user"
	PbAuth "FullBottle/auth/proto/auth"
)

var (
	c client.Client
)

func SetClient(client client.Client) {
	c = client
}

func GetUserSrvClient() PbUser.UserService {
	return PbUser.NewUserService("fullbottle.srv.user", c)
}

func GetAuthSrvClient() PbAuth.AuthService {
	return PbAuth.NewAuthService("fullbottle.srv.auth", c)
}
