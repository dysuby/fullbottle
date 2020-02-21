package main

import (
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/web"
	"github.com/vegchic/fullbottle/api/route"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/config"
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
