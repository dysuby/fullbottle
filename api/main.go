package main

import (
	"FullBottle/api/route"
	"FullBottle/common"
	"FullBottle/common/log"
	"FullBottle/config"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	service := web.NewService(
		web.Name(config.ApiName),
		web.Version("latest"),
	)

	if err := service.Init(); err != nil {
		log.Fatalf(err, "Service init failed")
	}

	common.SetClient(client.DefaultClient)

	service.Handle("/", route.GetGnRouter())

	// run service
	if err := service.Run(); err != nil {
		log.Fatalf(err, "Service running failed")
	}
}
