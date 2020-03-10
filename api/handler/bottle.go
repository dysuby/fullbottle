package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
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
		"msg":    "Success",
		"result": resp,
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
		req = &pbbottle.GetFolderInfoRequest{Ident: &pbbottle.GetFolderInfoRequest_Path_{Path: &pbbottle.GetFolderInfoRequest_Path{Relative: query.Path}}}
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
		"result": resp.Folder,
	})

}

func CreateFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	req := struct {
		Name     string `json:"name" binding:"required,max=100,min=1"`
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
		"msg":    "Success",
		"result": resp,
	})
}

func UpdateFolder(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FolderId int64  `json:"folder_id" binding:"required"`
		Name     string `json:"name" binding:"required,max=100,min=1"`
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
		FolderId int64 `json:"folder_id" binding:"required"`
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

func UpdateFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FileId   int64  `json:"file_id" binding:"required"`
		FolderId int64  `json:"folder_id" binding:"required"`
		Name     string `json:"name" binding:"required,max=100,min=1"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.UpdateFile(util.RpcContext(c), &pbbottle.UpdateFileRequest{FileId: body.FileId, OwnerId: uid,
		FolderId: body.FolderId, Name: body.Name})
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

func RemoveFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		FileId int64 `json:"file_id" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Body parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.RemoveFile(util.RpcContext(c), &pbbottle.RemoveFileRequest{OwnerId: uid, FileId: body.FileId})
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

func GetFolderParents(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		FolderId int64 `form:"folder_id" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid args",
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetEntryParents(util.RpcContext(c), &pbbottle.GetEntryParentsRequest{OwnerId: uid,
		EntryId: &pbbottle.GetEntryParentsRequest_FolderId{FolderId: query.FolderId}})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "Success",
		"result": resp.Parents,
	})
}
