package util

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func DetectContentType(c *gin.Context, f io.ReadSeeker) string {
	fileHeader := make([]byte, 512)
	if _, err := f.Read(fileHeader); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when operating file: " + err.Error(),
		})
		return ""
	}

	filetype := http.DetectContentType(fileHeader)

	if _, err := f.Seek(0, 0); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Error appears when operating file: " + err.Error(),
		})
		return ""
	}

	return filetype
}
