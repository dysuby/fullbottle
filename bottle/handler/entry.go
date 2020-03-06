package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
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

func (*EntryHandler) GetEntryParents(ctx context.Context, req *pb.GetEntryParentsRequest, resp *pb.GetEntryParentsResponse) error {
	ownerId := req.GetOwnerId()
	var folderId int64
	switch req.EntryId.(type) {
	case *pb.GetEntryParentsRequest_FileId:
		file, err := dao.GetFileById(ownerId, req.GetFileId())
		if err != nil {
			return err
		} else if file == nil {
			return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
		}
		folderId = file.FolderId
	case *pb.GetEntryParentsRequest_FolderId:
		folderId = req.GetFolderId()
	}

	resp.Parents = make([]*pb.GetEntryParentsResponseSimpleParent, 0)

	if folderId == dao.RootId {
		return nil
	}

	folder, err := dao.GetFolderById(ownerId, folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	resp.Parents = append(resp.Parents, &pb.GetEntryParentsResponseSimpleParent{FolderId:folder.ID, Name:folder.Name})	// add self to result
	for folder.ParentId != dao.RootId {
		folder, err = dao.GetFolderById(ownerId, folder.ParentId)
		if err != nil {
			return err
		}
		resp.Parents = append([]*pb.GetEntryParentsResponseSimpleParent{{FolderId:folder.ID, Name:folder.Name}}, resp.Parents...)
	}

	return nil
}
