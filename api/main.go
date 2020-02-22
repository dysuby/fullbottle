package main

import (
	"github.com/micro/go-micro/v2/client"
	gclient "github.com/micro/go-micro/v2/client/grpc"
	"github.com/micro/go-micro/v2/web"
	"github.com/vegchic/fullbottle/api/route"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
)

func cli() client.Client {
	return gclient.NewClient(
		client.WrapCall(common.ClientLogWrapper),
		gclient.MaxSendMsgSize(config.MaxMsgSendSize),
		gclient.MaxRecvMsgSize(config.MaxMsgRecvSize))
}


func main() {
	service := web.NewService(
		web.Name(config.ApiName),
		web.Version("latest"),
	)

	if err := service.Init(); err != nil {
		log.WithError(err).Fatalf("Service init failed")
	}

	common.SetClient(cli())

	service.Handle("/", route.GinRouter())

	// run service
	if err := service.Run(); err != nil {
		log.WithError(err).Fatalf("Service running failed")
	}
}
