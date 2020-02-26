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

	req := struct {
		FolderId int64  `json:"folder_id" bindings:"required"`
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
	_, err := bottleClient.UpdateFolder(util.RpcContext(c), &pbbottle.UpdateFolderRequest{OwnerId: uid, FolderId: req.FolderId,
		Name: req.Name, ParentId: req.ParentId})
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

	req := struct {
		FolderId int64 `json:"folder_id" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Body parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.RemoveFolder(util.RpcContext(c), &pbbottle.RemoveFolderRequest{OwnerId: uid, FolderId: req.FolderId})
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
