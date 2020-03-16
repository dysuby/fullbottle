package common

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/vegchic/fullbottle/config"

	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	pbshare "github.com/vegchic/fullbottle/share/proto/share"
	pbupload "github.com/vegchic/fullbottle/upload/proto/upload"
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

func BottleSrvClient() pbbottle.BottleService {
	return pbbottle.NewBottleService(config.BottleSrvName, c)
}

func UploadSrvClient() pbupload.UploadService {
	return pbupload.NewUploadService(config.UploadSrvName, c)
}

func ShareSrvClient() pbshare.ShareService {
	return pbshare.NewShareService(config.ShareSrvName, c)
}
