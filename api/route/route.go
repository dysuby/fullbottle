package route

import (
	"FullBottle/api/middleware"
	"FullBottle/api/handler"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(g *gin.Engine) {
	api := g.Group("/api")
	{
		api.POST("/register", handler.RegisterUser)
		api.POST("/login", handler.UserLogin)

		api.Use(middleware.LoginRequired)

		api.GET("/users/profile", handler.GetUser)
		api.PUT("/users/profile", handler.UpdateUser)
		api.POST("/users/profile/avatar", handler.UploadAvatar)
	}

}
