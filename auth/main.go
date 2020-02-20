package main

import (
	"github.com/vegchic/fullbottle/auth/handler"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
	"github.com/micro/go-micro/v2"

	auth "github.com/vegchic/fullbottle/auth/proto/auth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name(config.AuthSrvName),
		micro.Version("latest"),
		micro.WrapHandler(common.ServiceErrorRecovery),
		micro.WrapHandler(common.ServiceLogWrapper),
	)

	// Initialise service
	service.Init()

	// Register Handler
	if err := auth.RegisterAuthServiceHandler(service.Server(), new(handler.AuthHandler)); err != nil {
		log.Fatalf(err, "RegisterAuthServiceHandler failed")
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatalf(err, "Service running failed")
	}
}
