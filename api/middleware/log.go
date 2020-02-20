package middleware

import (
	"github.com/vegchic/fullbottle/common/log"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
	"time"
)

func ApiLogWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := time.Now()

		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf(err, "Failed to generate UUID")
		}
		log.WithFields(logrus.Fields{
			"uuid":     u.String(),
			"clientIP": c.ClientIP(),
		}).Infof("[%s] %s ", c.Request.Method, c.Request.URL.String())

		c.Next()

		log.WithFields(logrus.Fields{
			"uuid":     u.String(),
			"clientIP": c.ClientIP(),
			"cost":     (time.Now().UnixNano() - s.UnixNano()) / 1e6,
		}).Infof("%d ", c.Writer.Status())
	}
}
