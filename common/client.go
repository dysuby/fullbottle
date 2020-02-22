package common

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/vegchic/fullbottle/config"

	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	pbuser "github.com/vegchic/fullbottle/user/proto/user"
)

var (
	c client.Client
)

func SetClient(client client.Client) {
	c = client
}

func UserSrvClient() pbuser.UserService {
	return pbuser.NewUserService(config.UserSrvName, c)
}

func AuthSrvClient() pbauth.AuthService {
	return pbauth.NewAuthService(config.AuthSrvName, c)
}
