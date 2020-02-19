package main

import (
	"FullBottle/api/route"
	"FullBottle/common"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/util/log"
	"github.com/micro/go-micro/v2/web"
)


func main() {
	service := web.NewService(
		web.Name("fullbottle.web.api"),
		web.Version("latest"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	common.SetClient(client.DefaultClient)

	router := gin.Default()
	route.RegisterRoutes(router)

	service.Handle("/", router)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
