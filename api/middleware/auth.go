package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/api/util"
	pbauth "github.com/vegchic/fullbottle/auth/proto/auth"
	"github.com/vegchic/fullbottle/common"
	"net/http"
	"strings"
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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": e.Detail,
			})
			return
		}

		c.Set("cur_user_id", authResp.GetUserId())
		return
	}
}
