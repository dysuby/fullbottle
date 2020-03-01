package common

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/vegchic/fullbottle/config"

	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	pbuser "github.com/vegchic/fullbottle/user/proto/user"
	pbshare "github.com/vegchic/fullbottle/share/proto/share"
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

func BottleSrvClient() pbbottle.BottleService {
	return pbbottle.NewBottleService(config.BottleSrvName, c)
}

func ShareSrvClient() pbshare.ShareService {
	return pbshare.NewShareService(config.ShareSrvName, c)
}
