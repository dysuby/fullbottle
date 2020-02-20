package common

import (
	"github.com/vegchic/fullbottle/config"
	"github.com/micro/go-micro/v2/client"

	pbAuth "github.com/vegchic/fullbottle/auth/proto/auth"
	pbUser "github.com/vegchic/fullbottle/user/proto/user"
)

var (
	c client.Client
)

func SetClient(client client.Client) {
	c = client
}

func GetUserSrvClient() pbUser.UserService {
	return pbUser.NewUserService(config.UserSrvName, c)
}

func GetAuthSrvClient() pbAuth.AuthService {
	return pbAuth.NewAuthService(config.AuthSrvName, c)
}
