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
		Code      string  `json:"code" bindings:"required"`
		Exp       int64   `json:"exp" bindings:"required"`
		ParentId  int64	   `json:"parent_id" bindings:"required"`
		FileIds   []int64 `json:"file_ids" bindings:"required"`
		FolderIds []int64 `json:"folder_ids" bindings:"required"`
		IsPublic  bool 	  `json:"is_public" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.CreateShare(util.RpcContext(c), &pbshare.CreateShareRequest{
		SharerId:             uid,
		Code:                 body.Code,
		ParentId:             body.ParentId,
		FolderIds:            body.FolderIds,
		FileIds:              body.FileIds,
		Exp:                  body.Exp,
		IsPublic:             body.IsPublic,
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
		"result": resp,
	})
}


func UpdateShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token     string   `json:"token" bindings:"required"`
		Code      string  `json:"code" bindings:"required"`
		Exp       int64   `json:"exp" bindings:"required"`
		IsPublic  bool 	  `json:"is_public" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	_, err := shareClient.UpdateShare(util.RpcContext(c), &pbshare.UpdateShareRequest{
		Token:                   body.Token,
		SharerId:             uid,
		Code:                 body.Code,
		Exp:                  body.Exp,
		IsPublic:             body.IsPublic,
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

	body := struct {
		Token     string   `json:"token" bindings:"required"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	_, err := shareClient.CancelShare(util.RpcContext(c), &pbshare.CancelShareRequest{
		Token:                   body.Token,
		SharerId:             uid,
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
		Token        string   `uri:"token" bindings:"required"`
	}{}
	if err := c.ShouldBindUri(&param); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.ShareStatus(util.RpcContext(c), &pbshare.ShareStatusRequest{
		Token:                  param.Token,
		ViewerId:             	uid,
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
		"status": resp,
	})
}


func AccessShare(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Token        string   `json:"token"`
		Code         string   `json:"code"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.AccessShare(util.RpcContext(c), &pbshare.AccessShareRequest{
		Token:                   body.Token,
		ViewerId:             	uid,
		Code:					body.Code,
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
		"access_token": resp.AccessToken,
	})
}

func ShareInfo(c *gin.Context)  {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		AccessToken        string   `form:"access_token"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareInfo(util.RpcContext(c), &pbshare.GetShareInfoRequest{
		AccessToken:	query.AccessToken,
		ViewerId:             	uid,
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
		"info": resp,
	})
}

func ShareEntry(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	query := struct {
		AccessToken  string   `form:"access_token"`
		Path         string   `form:"path"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareFolder(util.RpcContext(c), &pbshare.GetShareFolderRequest{
		AccessToken:                   query.AccessToken,
		ViewerId:             	uid,
		Path:					query.Path,
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
		"folder": resp.Folder,
	})
}


func DownloadShareFile(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)
	body := struct {
		AccessToken        string   `json:"access_token"`
		FileId int64 `json:"file_id" bindings:"file_id"`
		Path         string   `json:"path"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	shareClient := common.ShareSrvClient()
	resp, err := shareClient.GetShareDownloadUrl(util.RpcContext(c), &pbshare.GetShareDownloadUrlRequest{
		AccessToken:          body.AccessToken,
		FileId:               body.FileId,
		Path:                 body.Path,
		ViewerId:             uid,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	util.DownloadProxy(c, resp.Url)
}
