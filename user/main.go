package main

import (
	"FullBottle/common"
	"FullBottle/user/handler"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/util/log"

	user "FullBottle/user/proto/user"
)

func main() {
	service := micro.NewService(
		micro.Name("fullbottle.srv.user"),
		micro.Version("latest"),
	)

	service.Init()

	common.SetClient(service.Client())

	if err := user.RegisterUserServiceHandler(service.Server(), new(handler.UserHandler)); err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
