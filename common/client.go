package common

import (
	"FullBottle/config"
	"github.com/micro/go-micro/v2/client"

	PbAuth "FullBottle/auth/proto/auth"
	PbUser "FullBottle/user/proto/user"
)

var (
	c client.Client
)

func SetClient(client client.Client) {
	c = client
}

func GetUserSrvClient() PbUser.UserService {
	return PbUser.NewUserService(config.UserSrvName, c)
}

func GetAuthSrvClient() PbAuth.AuthService {
	return PbAuth.NewAuthService(config.AuthSrvName, c)
}
