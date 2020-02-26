package util

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetIntVarWithAbort(c *gin.Context, name string) int {
	fn := []func(string) string{c.Param, c.Query}
	for _, f := range fn {
		if p := f(name); p == "" {
			continue
		} else if i, err := strconv.Atoi(p); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"msg": "Invalid " + name,
			})
			return 0
		} else {
			return i
		}
	}
	return 0
}
