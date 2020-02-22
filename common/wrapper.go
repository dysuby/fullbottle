package common

import (
	"context"
	"github.com/gofrs/uuid"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/metadata"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/server"
	"github.com/sirupsen/logrus"
	"github.com/vegchic/fullbottle/common/log"
	"time"
)

func ServiceLogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		s := time.Now()

		var reqId string
		if u, ok := metadata.Get(ctx, "uuid"); ok {
			reqId = u
		}
		if reqId == "" {
			u, err := uuid.NewV4()
			if err != nil {
				log.WithError(err).Warnf("Failed to generate UUID")
			} else {
				reqId = u.String()
			}
		}
		ctx = metadata.MergeContext(ctx, metadata.Metadata{
			"uuid": reqId,
		}, true)

		log.WithCtx(ctx).Infof("Recv rpc request: endpoint=%s", req.Endpoint())

		err := fn(ctx, req, resp)

		log.WithCtx(ctx).WithFields(logrus.Fields{
			"cost": (time.Now().UnixNano() - s.UnixNano()) / 1e6,
		}).Infof("Send rpc response err=%+v", err)

		return err
	}
}

func ServiceErrorRecovery(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		defer func() {
			if e := recover(); e != nil {
				log.Panic(e)
			}
		}()
		return fn(ctx, req, resp)
	}
}

func ClientLogWrapper(fn client.CallFunc) client.CallFunc {
	return func(ctx context.Context, node *registry.Node, req client.Request, resp interface{}, opts client.CallOptions) error {
		s := time.Now()

		log.WithCtx(ctx).Infof("Client rpc call send: endpoint=%s", req.Endpoint())

		err := fn(ctx, node, req, resp, opts)

		log.WithCtx(ctx).WithFields(logrus.Fields{
			"cost": (time.Now().UnixNano() - s.UnixNano()) / 1e6,
		}).Infof("Client rpc call finish: err=%v", err)

		return err
	}
}
