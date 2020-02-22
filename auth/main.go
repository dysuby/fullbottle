package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/vegchic/fullbottle/auth/handler"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"

	auth "github.com/vegchic/fullbottle/auth/proto/auth"
)

func options() []micro.Option {
	return []micro.Option{
		micro.Name(config.AuthSrvName),
		micro.Version("latest"),
		micro.WrapHandler(common.ServiceErrorRecovery),
		micro.WrapHandler(common.ServiceLogWrapper),
	}
}

func main() {
	// New Service
	service := micro.NewService(options()...)

	// Initialise service
	service.Init()

	// Register Handler
	if err := auth.RegisterAuthServiceHandler(service.Server(), new(handler.AuthHandler)); err != nil {
		log.WithError(err).Fatalf("RegisterAuthServiceHandler failed")
	}

	// Run service
	if err := service.Run(); err != nil {
		log.WithError(err).Fatalf("Service running failed")
	}
}
