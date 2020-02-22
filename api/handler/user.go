package handler

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	userdao "github.com/vegchic/fullbottle/user/dao"
	pbuser "github.com/vegchic/fullbottle/user/proto/user"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func GetUser(c *gin.Context) {
	u, _ := c.Get("CurrentUser")

	user := u.(*userdao.User)

	if user.Status == db.INVALID {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "Invalid user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"user": gin.H{
			"id":          user.ID,
			"status":      user.Status,
			"username":    user.Username,
			"email":       user.Email,
			"role":        user.Role,
			"avatar_fid":  user.AvatarFid,
			"create_time": user.CreateTime.Unix(),
		},
	})
}

func RegisterUser(c *gin.Context) {
	req := struct {
		Email    string `json:"email" binding:"required,email"`
		Username string `json:"username" binding:"required,max=24,min=4"`
		Password string `json:"password" binding:"required,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	_, err := client.CreateUser(util.RpcContext(c), &pbuser.CreateUserRequest{
		Email:    req.Email,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.EmailExisted {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "Email existed",
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": e.Detail,
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}

func UpdateUser(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	user := u.(*userdao.User)
	uid := user.ID

	req := struct {
		Username string `json:"username" binding:"required,max=24,min=4"`
		Password string `json:"password" binding:"required,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Arguments parse error: " + err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	_, err := client.ModifyUser(util.RpcContext(c), &pbuser.ModifyUserRequest{
		Uid:      uid,
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}

func UploadAvatar(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	user := u.(*userdao.User)
	uid := user.ID

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

	fbytes, err := ioutil.ReadAll(f)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when operating image: " + err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	req := &pbuser.UploadUserAvatarRequest{
		Uid:    uid,
		Avatar: fbytes,
	}

	if _, err = client.UploadUserAvatar(util.RpcContext(c), req); err != nil {
		e := errors.Parse(err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": e.Detail,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}

func GetUserAvatar(c *gin.Context) {
	q := c.Query("uid")
	uid, err := strconv.Atoi(q)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "invalid uid",
		})
		return
	}
	client := common.UserSrvClient()
	req := &pbuser.GetUserAvatarRequest{
		Uid: int64(uid),
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
	c.DataFromReader(http.StatusOK, int64(reader.Len()), "*", reader, map[string]string{})
	return
}

func UserLogin(c *gin.Context) {
	req := struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,max=18,min=6"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	client := common.UserSrvClient()
	resp, err := client.UserLogin(util.RpcContext(c), &pbuser.UserLoginRequest{Email: req.Email, Password: req.Password})
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.EmailExisted {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
				"msg": e.Detail,
			})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"msg": e.Detail,
			})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg":   "Success",
		"token": resp.Token,
	})
}
