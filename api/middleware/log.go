package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/vegchic/fullbottle/common/log"
	"time"
)

func ApiLogWrapper() gin.HandlerFunc {
	return func(c *gin.Context) {
		s := time.Now()

		log.WithCtx(c).Infof("[%s] %s ", c.Request.Method, c.Request.URL.String())

		c.Next()

		log.WithCtx(c).WithFields(logrus.Fields{
			"cost":     (time.Now().UnixNano() - s.UnixNano()) / 1e6,
		}).Infof("%d ", c.Writer.Status())
	}
}
