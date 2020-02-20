package common

import (
	"FullBottle/common/log"
	"context"
	"github.com/gofrs/uuid"
	"github.com/micro/go-micro/v2/server"
	"github.com/sirupsen/logrus"
	"time"
)

func ServiceLogWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		s := time.Now()

		u, err := uuid.NewV4()
		if err != nil {
			log.Fatalf(err, "Failed to generate UUID")
		}

		log.WithFields(logrus.Fields{
			"uuid": u.String(),
		}).Infof("Recv rpc request: endpoint=%s body=%s", req.Endpoint(), req.Body())

		ctx = context.WithValue(ctx, "uuid", u)
		err = fn(ctx, req, resp)

		log.WithFields(logrus.Fields{
			"uuid": u.String(),
			"cost": (time.Now().UnixNano() - s.UnixNano()) / 1e6,
		}).Infof("Send rpc response: %+v", resp)

		return err
	}
}

func ServiceErrorRecovery(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, resp interface{}) error {
		defer func() {
			if e := recover(); e!=nil {
				log.Panic(e)
			}
		}()
		return fn(ctx, req, resp)
	}
}
