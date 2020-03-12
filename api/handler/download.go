package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"net/http"
)

func CreateDownloadUrl(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)
	req := struct {
		FileId int64 `json:"file_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.CreateDownloadUrl(util.RpcContext(c), &pbbottle.CreateDownloadUrlRequest{FileId: req.FileId, OwnerId: uid, UserId: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "Success",
		"result": resp,
	})
}

func DownloadFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)
	param := struct {
		DownloadToken string `uri:"download_token" binding:"required"`
	}{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetWeedDownloadUrl(util.RpcContext(c), &pbbottle.GetWeedDownloadUrlRequest{DownloadToken: param.DownloadToken, UserId: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	util.DownloadProxy(c, resp.WeedUrl)
}

func GetImageThumbnail(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		FileId int64 `form:"file_id" binding:"required"`
		Height int32 `form:"mh"`
		Width  int32 `form:"mw"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid file_id",
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetImageThumbnail(util.RpcContext(c), &pbbottle.GetImageThumbnailRequest{
		FileId:  query.FileId,
		OwnerId: uid,
		UserId:  uid,
		Height:  query.Height,
		Width:   query.Width,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	b := resp.Thumbnail
	reader := bytes.NewReader(b)

	c.Header("Content-Type", "image/jpeg")
	c.DataFromReader(http.StatusOK, int64(reader.Len()), "image/jpeg", reader, map[string]string{})
}
