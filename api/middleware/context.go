package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/vegchic/fullbottle/common/log"
)

func WithContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		// generate request id
		u, err := uuid.NewV4()
		if err != nil {
			log.WithError(err).Fatalf("Failed to generate UUID")
		} else {
			c.Set("uuid", u.String())
		}

		// assign user metainfo
		c.Set("ip", c.ClientIP())
	}
}
