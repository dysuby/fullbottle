package main

import (
	"FullBottle/auth/handler"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"

	auth "FullBottle/auth/proto/auth"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("fullbottle.srv.auth"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	if err := auth.RegisterAuthServiceHandler(service.Server(), new(handler.AuthHandler)); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
