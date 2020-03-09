package util

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"mime/multipart"
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

func ReadFileBytes(c *gin.Context, fh *multipart.FileHeader) []byte {
	f, err := fh.Open()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Invalid file chunk",
		})
		return nil
	}
	defer f.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"msg": "Copy chunk failed",
		})
		return nil
	}
	return buf.Bytes()
}
