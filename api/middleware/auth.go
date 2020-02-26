package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/user/dao"
	pbuser "github.com/vegchic/fullbottle/user/proto/user"
	"net/http"
	"strings"
	"time"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate jwt token
		authClient := common.AuthSrvClient()
		authorization := c.GetHeader("authorization")
		if !strings.HasPrefix(authorization, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "Not authorized",
			})
			return
		}

		token := authorization[7:]
		authResp, err := authClient.ParseJwtToken(util.RpcContext(c), &pbauth.ParseJwtTokenRequest{Token: token})

		if err != nil {
			e := errors.Parse(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": e.Detail,
			})
			return
		}

		// get user info
		client := common.UserSrvClient()
		userResp, err := client.GetUserInfo(util.RpcContext(c), &pbuser.GetUserRequest{Uid: authResp.UserId})
		if err != nil {
			e := errors.Parse(err.Error())
			if e.Code == common.NotFoundError {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"msg": e.Detail,
				})
			} else {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"msg": e.Detail,
				})
			}
			return
		}

		ct := time.Unix(userResp.CreateTime, 0)
		ut := time.Unix(userResp.UpdateTime, 0)

		user := &dao.User{
			BasicModel: db.BasicModel{
				ID:         userResp.Uid,
				Status:     userResp.Status,
				CreateTime: &ct,
				UpdateTime: &ut,
			},
			Username:  userResp.Username,
			Password:  userResp.Password,
			Email:     userResp.Email,
			Role:      userResp.Role,
			AvatarFid: userResp.AvatarFid,
		}
		if userResp.DeleteTime != 0 {
			dt := time.Unix(userResp.DeleteTime, 0)
			user.DeleteTime = &dt
		}

		c.Set("CurrentUser", user)
		return
	}
}

func FolderAccessCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		var folderId int
		if f, ok := c.Get("folder_id"); ok {
			folderId = f.(int)
		}
		var userId int64
		if u, ok := c.Get("CurrentUser"); ok {
			userId = u.(*dao.User).ID
		}

		if folderId == 0 {
			return
		}

		authClient := common.AuthSrvClient()
		req := &pbauth.CheckFolderAccessRequest{FolderId:int64(folderId), UserId:userId}
		authResp, err := authClient.CheckFolderAccess(util.RpcContext(c), req)
		if err != nil {
			e := errors.Parse(err.Error())
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": e.Detail,
			})
			return
		}

		var need string
		if c.Request.Method == "GET" {
			need = config.ReadAction
		} else {
			need = config.WriteAction
		}
		for _, ac := range authResp.Actions {
			if ac == need {
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "No permission",
		})
	}
}
