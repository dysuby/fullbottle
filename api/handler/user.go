package handler

import (
	"fmt"
	"github.com/vegchic/fullbottle/api/utils"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/models"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"net/http"
	"strings"
	"time"

	pbUser "github.com/vegchic/fullbottle/user/proto/user"
)

func GetUser(c *gin.Context) {
	u, _ := c.Get("CurrentUser")

	user := u.(*models.User)

	if user.Status == models.INVALID {
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
			"avatar_url":  user.AvatarUrl,
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

	client := common.GetUserSrvClient()
	_, err := client.CreateUser(c, &pbUser.CreateUserRequest{
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
	user := u.(*models.User)
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

	client := common.GetUserSrvClient()
	_, err := client.ModifyUser(c, &pbUser.ModifyUserRequest{
		Id:       int64(uid),
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.UserNotFound {
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
		"msg": "Success",
	})
}

func UploadAvatar(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	user := u.(*models.User)
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
			"msg": "image size cannot exceed 2MB",
		})
		return
	}

	fileHeader := make([]byte, 512)
	if _, err := f.Read(fileHeader); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when operating image: " + err.Error(),
		})
		return
	}

	filetype := http.DetectContentType(fileHeader)
	if !strings.HasPrefix(filetype, "image") {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid image format",
		})
		return
	}

	if _, err := f.Seek(0, 0); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when operating image: " + err.Error(),
		})
		return
	}

	avatarName := fmt.Sprintf("%d-%d", uid, time.Now().Unix())
	fileInfo, err := utils.UploadFile(f, avatarName)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when uploading image: " + err.Error(),
		})
		return
	}

	client := common.GetUserSrvClient()
	_, err = client.ModifyUser(c, &pbUser.ModifyUserRequest{
		Id:        int64(uid),
		AvatarUrl: utils.JoinUrl(fileInfo.Url, fileInfo.Fid),
	})
	if err != nil {
		e := errors.Parse(err.Error())
		if e.Code == common.UserNotFound {
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
		"msg": "Success",
	})
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

	client := common.GetUserSrvClient()
	resp, err := client.UserLogin(c, &pbUser.UserLoginRequest{Email: req.Email, Password: req.Password})
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
