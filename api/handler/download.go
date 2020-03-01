package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"net/http"
)

func DownloadFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)
	req := struct {
		FileId int64 `json:"file_id" bindings:"file_id"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetDownloadUrl(util.RpcContext(c), &pbbottle.GetDownloadUrlRequest{FileId: req.FileId, OwnerId: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	util.DownloadProxy(c, resp.Url)
}
