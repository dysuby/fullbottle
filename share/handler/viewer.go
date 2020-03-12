package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/share/dao"
	pb "github.com/vegchic/fullbottle/share/proto/share"
	"github.com/vegchic/fullbottle/share/service"
)

type ViewerHandler struct{}

func (*ViewerHandler) GetShareInfo(ctx context.Context, req *pb.GetShareInfoRequest, resp *pb.GetShareInfoResponse) error {
	raw := req.GetAccessToken()
	viewerId := req.GetViewerId()
	at, err := service.ValidateAccessToken(raw, viewerId)
	if err != nil {
		return err
	}
	info, err := dao.GetShareById(at.SharerId, at.Id)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share not found", common.NotFoundError)
	} else if req.GetToken() != info.Token {
		return errors.New(config.ShareSrvName, "Share doesn't match this token", common.BadArgError)
	}

	metrics, err := dao.GetShareMetrics(info.ID)
	if err != nil {
		return err
	}

	resp.Id = info.ID
	resp.SharerId = info.SharerId
	if info.ExpireTime != nil {
		resp.Exp = info.ExpireTime.Unix()
	}
	for _, m := range metrics {
		if m.Action == dao.Download {
			resp.DownloadTimes = m.Times
		} else if m.Action == dao.View {
			resp.ViewTimes = m.Times
		}
	}

	return nil
}

func (*ViewerHandler) GetShareFolder(ctx context.Context, req *pb.GetShareFolderRequest, resp *pb.GetShareFolderResponse) error {
	viewerId := req.GetViewerId()
	raw := req.GetAccessToken()
	path := req.GetPath()

	at, err := service.ValidateAccessToken(raw, viewerId)
	if err != nil {
		return err
	}

	info, err := dao.GetShareById(at.SharerId, at.Id)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share not found", common.NotFoundError)
	} else if req.GetToken() != info.Token {
		return errors.New(config.ShareSrvName, "Share doesn't match this token", common.BadArgError)
	}

	folderResp, err := service.GetShareFolder(ctx, info, path)
	if err != nil {
		return err
	}

	resp.Folder = &pbbottle.FolderInfo{
		Folders: folderResp.Folder.Folders,
		Files:   folderResp.Folder.Files,
	}

	return nil
}

func (*ViewerHandler) GetShareDownloadUrl(ctx context.Context, req *pb.GetShareDownloadUrlRequest, resp *pb.GetShareDownloadUrlResponse) error {
	viewerId := req.GetViewerId()
	raw := req.GetAccessToken()

	at, err := service.ValidateAccessToken(raw, viewerId)
	if err != nil {
		return err
	}

	info, err := dao.GetShareById(at.SharerId, at.Id)
	if err != nil {
		return err
	} else if info == nil {
		return errors.New(config.ShareSrvName, "Share not found", common.NotFoundError)
	} else if req.GetToken() != info.Token {
		return errors.New(config.ShareSrvName, "Share doesn't match this token", common.BadArgError)
	}

	path := req.GetPath()
	fileId := req.GetFileId()

	folderResp, err := service.GetShareFolder(ctx, info, path)
	if err != nil {
		return err
	}

	v := false
	for _, f := range folderResp.Folder.GetFiles() {
		if f.Id == fileId {
			v = true
			break
		}
	}
	if !v {
		return errors.New(config.ShareSrvName, "File not found", common.NotFoundError)
	}

	bottleClient := common.BottleSrvClient()
	fileResp, err := bottleClient.CreateDownloadUrl(ctx, &pbbottle.CreateDownloadUrlRequest{OwnerId: info.SharerId, FileId: fileId, UserId: viewerId})
	if err != nil {
		return err
	}

	metric := &dao.ShareMetrics{
		ShareId:  info.ID,
		ViewerId: viewerId,
		Action:   dao.Download,
	}
	err = dao.CreateShareMetrics(metric)
	if err != nil {
		return err
	}

	resp.DownloadToken = fileResp.DownloadToken
	return nil
}
