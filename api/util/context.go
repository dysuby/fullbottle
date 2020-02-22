package util

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/micro/go-micro/v2/metadata"
)

func RpcContext(c *gin.Context) context.Context {
	ctx := context.Background()
	md := metadata.Metadata{
		"ip": c.ClientIP(),
	}

	if u, ok := c.Get("uuid"); ok {
		md["uuid"] = u.(string)
	}

	return metadata.NewContext(ctx, md)
}
