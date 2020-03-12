package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	"github.com/vegchic/fullbottle/common"
	pbshare "github.com/vegchic/fullbottle/share/proto/share"
	"net/http"
)

func CreateShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Code      string  `json:"code" binding:"max=6"`
		Exp       int64   `json:"exp"`
		ParentId  int64   `json:"parent_id" binding:"required"`
		FileIds   []int64 `json:"file_ids" binding:"required"`
		FolderIds []int64 `json:"folder_ids" binding:"required"`
		IsPublic  bool    `json:"is_public"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.CreateShare(util.RpcContext(c), &pbshare.CreateShareRequest{
		SharerId:  uid,
		Code:      body.Code,
		ParentId:  body.ParentId,
		FolderIds: body.FolderIds,
		FileIds:   body.FileIds,
		Exp:       body.Exp,
		IsPublic:  body.IsPublic,
	})
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

func UpdateShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token    string `json:"token" binding:"required"`
		Code     string `json:"code" binding:"max=6"`
		Exp      int64  `json:"exp"`
		IsPublic bool   `json:"is_public"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	_, err := shareClient.UpdateShare(util.RpcContext(c), &pbshare.UpdateShareRequest{
		Token:    body.Token,
		SharerId: uid,
		Code:     body.Code,
		Exp:      body.Exp,
		IsPublic: body.IsPublic,
	})
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

func CancelShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	param := struct {
		Token string `uri:"token" binding:"required"`
	}{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	_, err := shareClient.CancelShare(util.RpcContext(c), &pbshare.CancelShareRequest{
		Token:    param.Token,
		SharerId: uid,
	})
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

func RemoveShareEntries(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token     string  `json:"token" binding:"required"`
		FileIds   []int64 `json:"file_ids"`
		FolderIds []int64 `json:"folder_ids"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	_, err := shareClient.RemoveShareEntry(util.RpcContext(c), &pbshare.RemoveShareEntryRequest{
		Token:         body.Token,
		SharerId:      uid,
		RemoveFiles:   body.FileIds,
		RemoveFolders: body.FolderIds,
	})
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

func ShareStatus(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	param := struct {
		Token string `uri:"token" binding:"required"`
	}{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.ShareStatus(util.RpcContext(c), &pbshare.ShareStatusRequest{
		Token:    param.Token,
		ViewerId: uid,
	})
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

func AccessShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token string `json:"token" binding:"required"`
		Code  string `json:"code"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.AccessShare(util.RpcContext(c), &pbshare.AccessShareRequest{
		Token:    body.Token,
		ViewerId: uid,
		Code:     body.Code,
	})
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
			"access_token": resp.AccessToken,
			"exp":          resp.Exp,
		},
	})
}

func ShareInfo(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		AccessToken string `form:"access_token" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareInfo(util.RpcContext(c), &pbshare.GetShareInfoRequest{
		Token:       c.Param("token"),
		AccessToken: query.AccessToken,
		ViewerId:    uid,
	})
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

func ShareEntry(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		AccessToken string `form:"access_token"`
		Path        string `form:"path"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareFolder(util.RpcContext(c), &pbshare.GetShareFolderRequest{
		Token:       c.Param("token"),
		AccessToken: query.AccessToken,
		ViewerId:    uid,
		Path:        query.Path,
	})
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

func DownloadShareFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)
	body := struct {
		AccessToken string `json:"access_token" binding:"required"`
		FileId      int64  `json:"file_id" binding:"required"`
		Path        string `json:"path" binding:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareDownloadUrl(util.RpcContext(c), &pbshare.GetShareDownloadUrlRequest{
		Token:       c.Param("token"),
		AccessToken: body.AccessToken,
		FileId:      body.FileId,
		Path:        body.Path,
		ViewerId:    uid,
	})
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
