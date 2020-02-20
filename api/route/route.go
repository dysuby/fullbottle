package route

import (
	"github.com/vegchic/fullbottle/api/handler"
	"github.com/vegchic/fullbottle/api/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetGnRouter() *gin.Engine {
	router := gin.New()

	router.Use(middleware.ApiLogWrapper())
	router.Use(gin.Recovery())

	registerRoutes(router)

	return router
}

func registerRoutes(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.POST("/register", handler.RegisterUser)
		api.POST("/login", handler.UserLogin)

		api.Use(middleware.LoginRequired())

		api.GET("/users/profile", handler.GetUser)
		api.PUT("/users/profile", handler.UpdateUser)
		api.POST("/users/profile/avatar", handler.UploadAvatar)
	}

	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid api",
		})
	})

}
