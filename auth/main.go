package main

import (
	"FullBottle/auth/handler"
	"FullBottle/common"
	"FullBottle/common/log"
	"FullBottle/config"
	"github.com/micro/go-micro/v2"

	auth "FullBottle/auth/proto/auth"
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
