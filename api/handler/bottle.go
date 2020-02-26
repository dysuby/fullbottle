package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	userdao "github.com/vegchic/fullbottle/user/dao"
	"net/http"
)

func GetSpaceMeta(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	user := u.(*userdao.User)

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetBottleMetadata(util.RpcContext(c), &pbbottle.GetBottleMetadataRequest{Uid:user.ID})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"meta": resp,
	})
}

func GetFolder(c *gin.Context) {
	var folderId int
	if f, ok := c.Get("folder_id"); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Specify folder id first",
		})
		return
	} else {
		folderId = f.(int)
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.GetFolderInfo(util.RpcContext(c), &pbbottle.GetFolderInfoRequest{FolderId:int64(folderId)})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"folder": resp.Folder,
	})

}

func CreateFolder(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	user := u.(*userdao.User)

	req := struct {
		Name    string `json:"name" binding:"required,max=10,min=1"`
		ParentId int64 `json:"parentId" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	resp, err := bottleClient.CreateFolder(util.RpcContext(c), &pbbottle.CreateFolderRequest{Name:req.Name, ParentId:req.ParentId, OwnerId: user.ID})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"folderId": resp.GetFolderId(),
	})
}

func UpdateFolder(c *gin.Context) {
	var folderId int
	if f, ok := c.Get("folder_id"); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Specify folder id first",
		})
		return
	} else {
		folderId = f.(int)
	}

	req := struct {
		Name    string `json:"name" binding:"required,max=10,min=1"`
		ParentId int64 `json:"parentId" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.UpdateFolder(util.RpcContext(c), &pbbottle.UpdateFolderRequest{FolderId: int64(folderId),
		Name:req.Name, ParentId:req.ParentId})
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
	var folderId int
	if f, ok := c.Get("folder_id"); !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Specify folder id first",
		})
		return
	} else {
		folderId = f.(int)
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.RemoveFolder(util.RpcContext(c), &pbbottle.RemoveFolderRequest{FolderId: int64(folderId)})
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
