package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	gclient "github.com/micro/go-micro/v2/client/grpc"
	gserver "github.com/micro/go-micro/v2/server/grpc"
	"github.com/vegchic/fullbottle/upload/handler"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"

	pb "github.com/vegchic/fullbottle/upload/proto/upload"
)

func options() []micro.Option {
	return []micro.Option{
		micro.Server(gserver.NewServer( // need to pass first
			gserver.MaxMsgSize(config.MaxMsgSize),
		)),

		micro.Name(config.UploadSrvName),
		micro.Version("latest"),

		micro.Client(gclient.NewClient(
			client.WrapCall(common.ClientLogWrapper),
			gclient.MaxSendMsgSize(config.MaxMsgSendSize),
			gclient.MaxRecvMsgSize(config.MaxMsgRecvSize))),

		//micro.WrapHandler(common.ServiceErrorRecovery),
		micro.WrapHandler(common.ServiceLogWrapper),
	}
}

func main() {
	service := micro.NewService(options()...)

	service.Init()

	common.SetClient(service.Client())

	if err := pb.RegisterUploadHandler(service.Server(), new(handler.UploadHandler)); err != nil {
		log.WithError(err).Fatalf("RegisterBottleServiceHandler failed")
	}

	if err := service.Run(); err != nil {
		log.WithError(err).Fatalf("Service running failed")
	}
}
