package util

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func DownloadProxy(c *gin.Context, rawUrl string, filename string) {
	downloadUrl, err := url.Parse(rawUrl)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"msg": err.Error(),
		})
		return
	}
	weedReq, _ := http.NewRequest("GET", downloadUrl.String(), bytes.NewReader([]byte{}))
	proxy := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			req.Header = c.Request.Header
		},
		ModifyResponse: func(r *http.Response) error {
			// make it download
			r.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))
			return nil
		},
	}

	proxy.ServeHTTP(c.Writer, weedReq)
}
