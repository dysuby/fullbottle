package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"io"
	"mime/multipart"
	"net/http"
)

func GetSpaceMeta(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetBottleMetadata(util.RpcContext(c), &pbbottle.GetBottleMetadataRequest{Uid: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":  "Success",
		"meta": resp,
	})
}

func GetFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		Path     string `form:"path"`
		FolderId int64  `form:"folder_id"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid args",
		})
		return
	}

	var req *pbbottle.GetFolderInfoRequest
	if query.FolderId != 0 {
		req = &pbbottle.GetFolderInfoRequest{Ident: &pbbottle.GetFolderInfoRequest_FolderId{FolderId: query.FolderId}}
	} else {
		req = &pbbottle.GetFolderInfoRequest{Ident: &pbbottle.GetFolderInfoRequest_Path{Path: query.Path}}
	}
	req.OwnerId = uid
	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetFolderInfo(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "Success",
		"folder": resp.Folder,
	})

}

func CreateFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	req := struct {
		Name     string `json:"name" binding:"required,max=10,min=1"`
		ParentId int64  `json:"parent_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.CreateFolder(util.RpcContext(c), &pbbottle.CreateFolderRequest{Name: req.Name, ParentId: req.ParentId, OwnerId: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":      "Success",
		"folderId": resp.GetFolderId(),
	})
}

func UpdateFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FolderId int64  `json:"folder_id" bindings:"required"`
		Name     string `json:"name" binding:"required,max=10,min=1"`
		ParentId int64  `json:"parent_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.UpdateFolder(util.RpcContext(c), &pbbottle.UpdateFolderRequest{OwnerId: uid, FolderId: body.FolderId,
		Name: body.Name, ParentId: body.ParentId})
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

func RemoveFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FolderId int64 `json:"folder_id" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Body parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.RemoveFolder(util.RpcContext(c), &pbbottle.RemoveFolderRequest{OwnerId: uid, FolderId: body.FolderId})
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

	req := &pbbottle.GenerateUploadTokenRequest{
		OwnerId:  uid,
		Filename: body.Filename,
		FolderId: body.FolderId,
		Hash:     body.Hash,
		Size:     body.Size,
		Mime:     body.Mime,
	}
	resp, err := bottleClient.GenerateUploadToken(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Success",
		"result": gin.H{
			"token": resp.Token,
			"need_upload": resp.NeedUpload,
			"uploaded": resp.Uploaded,
		},
	})
}

func UploadFile(c *gin.Context) {
	body := struct {
		FilePart *multipart.FileHeader `form:"file_part"`
		Token    string                `form:"token" bindings:"required"`
		Offset   int64                 `form:"offset" bindings:"required"`
	}{}
	if err := c.ShouldBind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid token",
		})
		return
	}
	var b []byte
	if body.FilePart != nil {
		f, err := body.FilePart.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "Invalid file",
			})
			return
		}
		defer f.Close()

		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "Invalid file",
			})
			return
		}
		b = buf.Bytes()
	}

	bottleClient := common.BottleSrvClient()
	req := &pbbottle.UploadFileRequest{
		Token:  body.Token,
		Offset: body.Offset,
		Raw:    b,
	}
	resp, err := bottleClient.UploadFile(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "Success",
		"result": gin.H{
			"status": resp.Status,
			"uploaded": resp.Uploaded,
		},
	})
}
