package util

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func DownloadProxy(c *gin.Context, rawUrl string) {
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

		},
		ModifyResponse: func(r *http.Response) error {
			// make it download
			cd := r.Header.Get("Content-Disposition")
			r.Header.Set("Content-Disposition", strings.Replace(cd, "inline;", "attachment;", 1))
			return nil
		},
	}

	proxy.ServeHTTP(c.Writer, weedReq)
}
