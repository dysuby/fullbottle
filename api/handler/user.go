package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	pbuser "github.com/vegchic/fullbottle/user/proto/user"
	"net/http"
	"strings"
)

func GetUser(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	// get user info
	client := common.UserSrvClient()
	userResp, err := client.GetUserInfo(util.RpcContext(c), &pbuser.GetUserRequest{Uid: uid})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})
		return
	}

	if userResp.Status == db.Invalid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":    "Success",
		"result": userResp,
	})
}

func RegisterUser(c *gin.Context) {
	body := struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,max=24,min=4"`
		Password string `json:"password" binding:"required,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	_, err := client.CreateUser(util.RpcContext(c), &pbuser.CreateUserRequest{
		Email:    body.Email,
		Username: body.Username,
		Password: body.Password,
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

func UpdateUser(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	body := struct {
		Username string `json:"username" binding:"max=24,min=4"`
		Password string `json:"password" binding:"omitempty,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	_, err := client.ModifyUser(util.RpcContext(c), &pbuser.ModifyUserRequest{
		Uid:      uid,
		Username: body.Username,
		Password: body.Password,
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

func UploadAvatar(c *gin.Context) {
	u, _ := c.Get("cur_user_id")
	uid := u.(int64)

	f, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Cannot read image from form",
		})
		return
	}
	if header.Size > config.AvatarMaxSize {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "image size cannot exceed 1MB",
		})
		return
	}

	filetype := util.DetectContentType(c, f)
	if c.IsAborted() {
		return
	}

	if !strings.HasPrefix(filetype, "image") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid image format",
		})
		return
	}

	fbytes := util.ReadFileBytes(c, header)
	if c.IsAborted() {
		return
	}

	client := common.UserSrvClient()
	req := &pbuser.UploadUserAvatarRequest{
		Uid:    uid,
		Avatar: fbytes,
	}

	if _, err = client.UploadUserAvatar(util.RpcContext(c), req); err != nil {
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

func GetUserAvatar(c *gin.Context) {
	query := struct {
		Uid int64 `form:"uid" binding:"required"`
	}{}
	if err := c.ShouldBindQuery(&query); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid uid",
		})
		return
	}

	client := common.UserSrvClient()
	req := &pbuser.GetUserAvatarRequest{
		Uid: query.Uid,
	}
	resp, err := client.GetUserAvatar(util.RpcContext(c), req)
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.EmptyAvatarError {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg": "avatar not found: " + e.Detail,
		})
		return
	}

	b := resp.Avatar
	reader := bytes.NewReader(b)

	c.Header("Content-Type", resp.ContentType)
	c.DataFromReader(http.StatusOK, int64(reader.Len()), resp.ContentType, reader, map[string]string{})
}

func UserLogin(c *gin.Context) {
	body := struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	resp, err := client.UserLogin(util.RpcContext(c), &pbuser.UserLoginRequest{Email: body.Email, Password: body.Password})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": e.Detail,
		})

		return
	}

	c.SetCookie("token", resp.Token, int(resp.Expire), "/", config.C().Server.Ip, false, true)

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"result": gin.H{
			"token":  resp.Token,
			"expire": resp.Expire,
			"uid":    resp.Uid,
		},
	})
}
