package middleware

import (
	PbAuth "FullBottle/auth/proto/auth"
	"FullBottle/common"
	PbUser "FullBottle/user/proto/user"
	"github.com/micro/go-micro/v2/util/log"
	"strings"

	"github.com/gin-gonic/gin"
	"net/http"
)

func LoginRequired(c *gin.Context) {
	// validate jwt token
	authClient := common.GetAuthSrvClient()
	authorization := c.GetHeader("authorization")
	if !strings.HasPrefix(authorization, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"msg": "Not authorized",
		})
		return
	}

	token := authorization[7:]
	log.Info(token)
	authResp, err := authClient.ParseJwtToken(c, &PbAuth.ParseJwtTokenRequest{Token:token})
	log.Info(err, authResp)
	if err != nil || authResp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Cannot parse jwt token",
		})
		return
	}

	// get user info
	client := common.GetUserSrvClient()
	userResp, err := client.GetUserInfo(c, &PbUser.GetUserRequest{Id: authResp.UserId})
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}

	if userResp.Code == common.UserNotFound {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"msg": userResp.Msg,
		})
		return
	} else if userResp.Code != common.Success {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": userResp.Msg,
		})
		return
	}

	c.Set("CurrentUser", *userResp.GetInfo())
}
