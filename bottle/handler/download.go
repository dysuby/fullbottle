package handler

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/kv"
	"github.com/vegchic/fullbottle/config"
	"github.com/vegchic/fullbottle/util"
	"github.com/vegchic/fullbottle/weed"
	"time"
)

const DownloadTokenKey = "download:token=%s;user_id=%d"

type DownloadHandler struct{}

func (*DownloadHandler) CreateDownloadUrl(ctx context.Context, req *pb.CreateDownloadUrlRequest, resp *pb.CreateDownloadUrlResponse) error {
	fileId := req.GetFileId()
	ownerId := req.GetOwnerId()
	userId := req.GetUserId()

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

	token := util.GenToken(20)
	if err := kv.C().Set(fmt.Sprintf(DownloadTokenKey, token, userId), downloadUrl.String(), 5*time.Minute).Err(); err != nil {
		return common.NewRedisError(err)
	}

	resp.DownloadToken = token
	return nil
}

func (*DownloadHandler) GetWeedDownloadUrl(ctx context.Context, req *pb.GetWeedDownloadUrlRequest, resp *pb.GetWeedDownloadUrlResponse) error {
	token := req.GetDownloadToken()
	userId := req.GetUserId()

	res, err := kv.C().Get(fmt.Sprintf(DownloadTokenKey, token, userId)).Result()
	if err != nil {
		return common.NewRedisError(err)
	}

	resp.WeedUrl = res
	return nil
}
