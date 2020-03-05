package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
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

	var expire time.Time
	if exp != 0 {
		expire = time.Unix(exp, 0)
		if expire.Before(time.Now()) {
			return errors.New(config.ShareSrvName, "Invalid expire: "+expire.String(), common.BadArgError)
		}
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
		break
	}

	info := &dao.ShareInfo{
		SharerId:       sharerId,
		Token:          token,
		ParentFolderId: parentId,
		ShareRefs:      refs,
	}
	if !expire.IsZero() {
		info.ExpireTime = &expire
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
	token := req.GetToken()
	sharerId := req.GetSharerId()

	code := req.GetCode()
	exp := req.GetExp()

	info, err := dao.GetShareByToken(token)
	if err != nil {
		return err
	} else if info == nil || info.SharerId != sharerId {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	var expire time.Time
	if exp != 0 {
		expire = time.Unix(exp, 0)
		if expire.Before(time.Now()) {
			return errors.New(config.ShareSrvName, "Invalid expire: "+expire.String(), common.BadArgError)
		}
	}
	if !expire.IsZero() {
		info.ExpireTime = &expire
	}

	if !req.GetIsPublic() {
		info.Privacy = dao.Private
		info.Code = util.TokenMd5(code)
	} else {
		info.Privacy = dao.Public
	}

	return dao.UpdateShare(info)
}

func (*SharerHandler) CancelShare(ctx context.Context, req *pb.CancelShareRequest, resp *pb.CancelShareResponse) error {
	token := req.GetToken()
	sharerId := req.GetSharerId()

	info, err := dao.GetShareByToken(token)
	if err != nil {
		return err
	} else if info == nil || info.SharerId != sharerId {
		return errors.New(config.ShareSrvName, "Share info not found", common.NotFoundError)
	}

	return dao.CancelShare(info, db.Canceled)
}
