package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/weed"
)

type DownloadHandler struct{}

func (*DownloadHandler) GetDownloadUrl(ctx context.Context, req *pb.GetDownloadUrlRequest, resp *pb.GetDownloadUrlResponse) error {
	fileId := req.GetFileId()
	ownerId := req.GetOwnerId()

	file, err := dao.GetFileById(ownerId, fileId)
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	fid := file.Metadata.Fid
	downloadUrl, err := weed.GetDownloadUrl(fid)
	if err != nil {
		return err
	}

	resp.Url = downloadUrl.String()
	return nil
}
