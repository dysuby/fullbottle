package handler

import (
	"FullBottle/api/utils"
	"FullBottle/common"
	"FullBottle/models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"

	PbUser "FullBottle/user/proto/user"
)

const AvatarMaxSize = 2 * (1 << 20)  // 2mb

func GetUser(c *gin.Context) {
	u, _ := c.Get("CurrentUser")

	info := u.(PbUser.UserInfo)

	if info.Status == models.INVALID {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": "Invalid user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"user": gin.H{
			"id":          info.Id,
			"status":      info.Status,
			"username":    info.Username,
			"email":       info.Email,
			"avatar_uri":  info.AvatarUri,
			"create_time": info.CreateTime,
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
			"msg": err.Error(),
		})
		return
	}

	client := common.GetUserSrvClient()
	resp, err := client.CreateUser(c, &PbUser.CreateUserRequest{
		Var: &PbUser.UserVar{
			Email:    req.Email,
			Username: req.Username,
			Password: req.Password,
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if resp.Code == common.EmailExisted {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Email existed",
		})
		return
	} else if resp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": resp.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}

func UpdateUser(c *gin.Context) {
	u, _ := c.Get("CurrentUser")
	info := u.(PbUser.UserInfo)
	uid := info.Id

	req := struct {
		Username string `json:"username" binding:"required,max=24,min=4"`
		Password string `json:"password" binding:"required,max=18,min:6"`
	}{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
		return
	}

	client := common.GetUserSrvClient()
	resp, err := client.ModifyUser(c, &PbUser.ModifyUserRequest{
		Id: int64(uid),
		Var: &PbUser.UserVar{
			Username: req.Username,
			Password: req.Password,
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if resp.Code == common.UserNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg": resp.Msg,
		})
		return
	} else if resp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": resp.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
	})
}

func UploadAvatar(c *gin.Context)  {
	u, _ := c.Get("CurrentUser")
	info := u.(PbUser.UserInfo)
	uid := info.Id

	f, header, err := c.Request.FormFile("avatar")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Cannot read image from form",
		})
		return
	}
	if header.Size > AvatarMaxSize {
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
	_, err = utils.UploadFile(f, avatarName, "__avatar__")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when uploading image: " + err.Error(),
		})
		return
	}

	client := common.GetUserSrvClient()
	resp, err := client.ModifyUser(c, &PbUser.ModifyUserRequest{
		Id: int64(uid),
		Var: &PbUser.UserVar{
			AvatarUri: utils.GenFilePath("__avatar__", avatarName),
		},
	})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if resp.Code == common.UserNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg": resp.Msg,
		})
		return
	} else if resp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": resp.Msg,
		})
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
	resp, err := client.UserLogin(c, &PbUser.UserLoginRequest{Email: req.Email, Password:req.Password})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if resp.Code == common.UserNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg": resp.Msg,
		})
		return
	} else if resp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": resp.Msg,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"msg": "Success",
		"token": resp.Token,
	})
}