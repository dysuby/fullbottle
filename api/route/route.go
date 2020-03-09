package route

import (
	"github.com/gin-gonic/gin"
	"github.com/vegchic/fullbottle/api/handler"
	"github.com/vegchic/fullbottle/api/middleware"
	"net/http"
)

func GinRouter() *gin.Engine {
	router := gin.New()

	router.Use(gin.Recovery())
	router.Use(middleware.WithContext())
	router.Use(middleware.ApiLogWrapper())

	registerV1Routes(router)

	return router
}

func registerV1Routes(g *gin.Engine) {
	api := g.Group("/v1")
	{
		api.POST("/register", handler.RegisterUser)
		api.POST("/login", handler.UserLogin)

		api.GET("/users/avatar", handler.GetUserAvatar) // no permission asked

		api.Use(middleware.LoginRequired())

		// user
		api.GET("/users/profile", handler.GetUser)
		api.PUT("/users/profile", handler.UpdateUser)
		api.POST("/users/avatar", handler.UploadAvatar)

		// space
		api.GET("/space/meta", handler.GetSpaceMeta)
		api.GET("/space/folders", handler.GetFolder) // ?folder_id=&path=
		api.POST("/space/folders", handler.CreateFolder)
		api.PUT("/space/folders", handler.UpdateFolder)
		api.DELETE("/space/folders", handler.RemoveFolder)

		api.GET("/space/folders/parents", handler.GetFolderParents)

		api.PUT("/space/files", handler.UpdateFile)
		api.DELETE("/space/files", handler.RemoveFile)

		api.POST("/space/upload/token", handler.GetUploadToken) // ask for token
		api.POST("/space/upload/file", handler.UploadFile)      // upload file
		api.GET("/space/upload/file", handler.GetUploadedFileChunks)
		api.DELETE("/space/upload/file", handler.CancelFileUpload)

		api.POST("/space/download/file", handler.CreateDownloadUrl)

		// share
		api.POST("/share", handler.CreateShare)          // create share
		api.PUT("/share/:token", handler.UpdateShare)    // modify share
		api.DELETE("/share/:token", handler.CancelShare) // cancel share

		api.GET("/share/:token/status", handler.ShareStatus)  // get share status
		api.POST("/share/:token/access", handler.AccessShare) // try to access share
		api.GET("/share/:token/info", handler.ShareInfo)
		api.GET("/share/:token", handler.ShareEntry)                  // get share entries
		api.POST("/share/:token/download", handler.DownloadShareFile) // download

		// raw download
		api.GET("/download/file/:download_token", handler.DownloadFile)
	}

	g.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid api",
		})
	})

}
