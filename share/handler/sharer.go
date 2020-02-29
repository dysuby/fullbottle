package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/share/dao"
	pb "github.com/vegchic/fullbottle/share/proto/share"
	"github.com/vegchic/fullbottle/share/service"
	"github.com/vegchic/fullbottle/util"
	"time"
)

type SharerHandler struct{}

func (*SharerHandler) CreateShare(ctx context.Context, req *pb.CreateShareRequest, resp *pb.CreateShareResponse) error {
	sharerId := req.GetSharerId()
	code := req.GetCode()
	exp := req.GetExp()
	parentId := req.GetParentId()
	folderIds := req.GetFolderIds()
	fileIds := req.GetFileIds()

	expire := time.Unix(exp, 0)
	if !expire.IsZero() && expire.Before(time.Now()) {
		return errors.New(config.ShareSrvName, "Invalid expire", common.BadArgError)
	}

	refs, err := service.ValidateEntries(ctx, sharerId, parentId, folderIds, fileIds)
	if err != nil {
		return err
	}

	token := util.GenToken(10)
	for true {
		if o, err := dao.GetShareByToken(token); err != nil {
			return err
		} else if o != nil {
			token = util.GenToken(10)
		}
	}

	info := &dao.ShareInfo{
		SharerId:       sharerId,
		Token:          token,
		ExpireTime:     &expire,
		ParentFolderId: parentId,
		ShareRefs:      refs,
	}
	if !req.GetIsPublic() {
		info.Privacy = dao.Private
		info.Code = util.TokenMd5(code)
	}
	err = dao.InitShare(info)
	if err != nil {
		return err
	}

	resp.Id = info.ID
	resp.Token = info.Token
	return nil
}

func (*SharerHandler) UpdateShare(ctx context.Context, req *pb.UpdateShareRequest, resp *pb.UpdateShareResponse) error {
	id := req.GetId()
	sharerId := req.GetSharerId()

	code := req.GetCode()
	exp := req.GetExp()

	info, err := dao.GetShareById(sharerId, id)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	expire := time.Unix(exp, 0)
	if !expire.IsZero() && expire.Before(time.Now()) {
		return errors.New(config.ShareSrvName, "Invalid expire", common.BadArgError)
	}

	info.ExpireTime = &expire
	if !req.GetIsPublic() {
		info.Privacy = dao.Private
		info.Code = util.TokenMd5(code)
	}

	return dao.UpdateShare(info)
}

func (*SharerHandler) CancelShare(ctx context.Context, req *pb.CancelShareRequest, resp *pb.CancelShareResponse) error {
	id := req.GetId()
	sharerId := req.GetSharerId()

	info, err := dao.GetShareById(sharerId, id)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	return dao.CancelShare(info)
}
