package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/util"
	"net/http"
)

func LoginRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		// validate jwt token
		token, err := c.Cookie("token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "Not authorized",
			})
			return
		}

		ip, _ := c.Get("ip")
		claims, err := util.ParseJwtToken(token, ip.(string))

		if err != nil {
			e := errors.Parse(err.Error())
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": e.Detail,
			})
			return
		}

		c.Set("cur_user_id", claims.Uid)
		return
	}
}
