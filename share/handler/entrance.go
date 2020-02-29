package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/share/dao"
	pb "github.com/vegchic/fullbottle/share/proto/share"
	"github.com/vegchic/fullbottle/share/service"
	"github.com/vegchic/fullbottle/util"
	"time"
)

type EntranceHandler struct{}

func (*EntranceHandler) ShareStatus(ctx context.Context, req *pb.ShareStatusRequest, resp *pb.ShareStatusResponse) error {
	token := req.GetToken()

	info, err := dao.GetShareByToken(token)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	resp.Status = info.Status
	resp.OwnerId = info.SharerId
	resp.IsPublic = info.Privacy == dao.Public

	return nil
}

func (*EntranceHandler) AccessShare(ctx context.Context, req *pb.AccessShareRequest, resp *pb.AccessShareResponse) error {
	viewerId := req.GetViewerId()
	token := req.GetToken()
	code := req.GetCode()

	info, err := dao.GetShareByToken(token)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	if util.TokenMd5(code) != info.Code {
		return errors.New(config.ShareSrvName, "Error code", common.BadArgError)
	}

	at := service.NewAccessToken(info.SharerId, info.ID, viewerId)
	if err := kv.Set(at.Token, at, 24*time.Hour); err != nil {
		return err
	}

	metric := &dao.ShareMetrics{ViewerId: viewerId, ShareId: info.ID, Action: dao.View}
	if err = dao.CreateShareMetrics(metric); err != nil {
		return err
	}

	resp.AccessToken = at.Token
	return nil
}
