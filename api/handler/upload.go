package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	pbupload "github.com/vegchic/fullbottle/upload/proto/upload"
	"mime/multipart"
	"net/http"
	"strings"
)

func GetUploadToken(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FolderId int64  `json:"folder_id" bindings:"required"`
		Filename string `json:"filename" bindings:"required,min=1,max=100"`
		Mime     string `json:"mime" bindings:"required"`
		Hash     string `json:"hash" bindings:"required"`
		Size     int64  `json:"size" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if !strings.Contains(body.Mime, "/") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid mime type",
		})
		return
	}

	bottleClient := common.BottleSrvClient()

	bottleMeta, err := bottleClient.GetBottleMetadata(util.RpcContext(c), &pbbottle.GetBottleMetadataRequest{Uid: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}
	if bottleMeta.Remain < body.Size {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Remain space not enough",
		})
		return
	}

	uploadClient := common.UploadSrvClient()
	req := &pbupload.GenerateUploadTokenRequest{
		OwnerId:  uid,
		Filename: body.Filename,
		FolderId: body.FolderId,
		Hash:     body.Hash,
		Size:     body.Size,
		Mime:     body.Mime,
	}
	resp, err := uploadClient.GenerateUploadToken(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"result": gin.H{
			"token":       resp.Token,
			"need_upload": resp.NeedUpload,
			"uploaded":    resp.Uploaded,
		},
	})
}

func UploadFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		File      *multipart.FileHeader `form:"file"`
		Token     string                `form:"token" bindings:"required"`
		Offset    int64                 `form:"offset" bindings:"required"`
		ChunkHash string                `form:"chunk_hash" bindings:"required"`
	}{}
	if err := c.ShouldBind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid token",
		})
		return
	}

	uploadClient := common.UploadSrvClient()
	uploadedResp, err := uploadClient.GetFileUploadedChunks(util.RpcContext(c), &pbupload.GetFileUploadedChunksRequest{Token: body.Token, OwnerId:uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}
	for _, offset := range uploadedResp.Uploaded {
		if offset == body.Offset {
			c.JSON(http.StatusCreated, gin.H{
				"msg": "The chunk has been uploaded",
				"result": gin.H{
					"uploaded": uploadedResp.Uploaded,
					"need_upload": uploadedResp.NeedUpload,
				},
			})
			return
		}
	}

	var b []byte
	if body.File != nil {
		b = util.ReadFileBytes(c, body.File)
		if c.IsAborted() {
			return
		}
	}

	req := &pbupload.UploadFileRequest{
		Token:  body.Token,
		OwnerId:uid,
		Offset: body.Offset,
		Raw:    b,
		ChunkHash:body.ChunkHash,
	}
	resp, err := uploadClient.UploadFile(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.FileFailError {
			// todo use a better status code
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg": e.Detail,
			})
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"result": gin.H{
			"status":   resp.Status,
			"uploaded": resp.Uploaded,
		},
	})
}

func GetUploadedFileChunks(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		Token string `form:"token" bindings:"required"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid token",
		})
		return
	}
	uploadClient := common.UploadSrvClient()
	resp, err := uploadClient.GetFileUploadedChunks(util.RpcContext(c), &pbupload.GetFileUploadedChunksRequest{Token: query.Token, OwnerId:uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"result": gin.H{
			"uploaded": resp.Uploaded,
			"need_upload": resp.NeedUpload,
		},
	})
}

func CancelFileUpload(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token    string                `form:"token" bindings:"required"`
	}{}
	if err := c.ShouldBind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid token",
		})
		return
	}
	uploadClient := common.UploadSrvClient()
	_, err := uploadClient.CancelFileUpload(util.RpcContext(c), &pbupload.CancelFileUploadRequest{Token:body.Token, OwnerId:uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}
