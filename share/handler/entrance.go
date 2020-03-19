package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
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

	info, err := dao.GetShareByToken(token, false)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	if info.ExpireTime != nil && info.ExpireTime.Before(time.Now()) {
		if err := dao.CancelShare(info, db.Expired); err != nil {
			return err
		}
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

	info, err := dao.GetShareByToken(token, true)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	var maxExp int64
	if info.ExpireTime != nil && info.ExpireTime.Before(time.Now()) {
		if err := dao.CancelShare(info, db.Expired); err != nil {
			return err
		}
		return errors.New(config.ShareSrvName, "Share expired", common.NotFoundError)
	} else if info.ExpireTime != nil {
		maxExp = info.ExpireTime.Unix()
	}

	// permission check
	if info.SharerId != viewerId && info.Privacy != dao.Public && util.Md5(code) != info.Code {
		return errors.New(config.ShareSrvName, "Access request denied", common.BadArgError)
	}

	at := service.NewAccessToken(info.ID, info.SharerId, viewerId)
	if err := kv.SetM(at.Token, at, 24*time.Hour); err != nil {
		return err
	}

	metric := &dao.ShareMetrics{ViewerId: viewerId, ShareId: info.ID, Action: dao.View}
	if err = dao.CreateShareMetrics(metric); err != nil {
		return err
	}

	resp.AccessToken = at.Token

	exp := time.Now().Add(24 * time.Hour).Unix()
	if maxExp != 0 && maxExp < exp {
		resp.Exp = maxExp
	} else {
		resp.Exp = exp
	}
	return nil
}
