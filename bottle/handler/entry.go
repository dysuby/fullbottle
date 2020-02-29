package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/bottle/service"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/config"
)

type EntryHandler struct{}

func (*EntryHandler) ValidateEntry(ctx context.Context, req *pb.ValidateEntryRequest, resp *pb.ValidateEntryResponse) error {
	ownerId := req.GetOwnerId()
	parentId := req.GetParentId()
	folderIds := req.GetFolderIds()
	fileIds := req.GetFileIds()

	folders, files, err := service.GetSubEntry(ownerId, parentId, folderIds, fileIds)
	if err != nil {
		return err
	}
	if folderIds != nil && len(folderIds) != len(folders) {
		return errors.New(config.BottleSrvName, "Some folder not found", common.NotFoundError)
	}
	if fileIds != nil && len(fileIds) != len(files) {
		return errors.New(config.BottleSrvName, "Some file not found", common.NotFoundError)
	}

	return nil
}
