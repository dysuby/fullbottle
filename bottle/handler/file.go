package handler

import (
	"context"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	"time"
)

type FileHandler struct{}

func (f *FileHandler) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest, resp *pb.GetFileInfoResponse) error {
	file, err := dao.GetFileById(req.GetOwnerId(), req.GetFileId())

	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	resp.File = &pb.FileInfo{
		Id:         file.ID,
		FileId:     file.FileId,
		Name:       file.Name,
		Size:       file.Metadata.Size,
		Hash:       file.Metadata.Hash,
		FolderId:   file.FolderId,
		OwnerId:    file.OwnerId,
		CreateTime: file.CreateTime.Unix(),
		UpdateTime: file.UpdateTime.Unix(),
	}
	return nil
}

func (f *FileHandler) UpdateFile(ctx context.Context, req *pb.UpdateFileRequest, resp *pb.UpdateFileResponse) error {
	name := req.GetName()
	folderId := req.GetFolderId()
	fileId := req.GetFileId()
	ownerId := req.GetOwnerId()

	file, err := dao.GetFileById(ownerId, fileId)
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}
	folder, err := dao.GetFolderById(ownerId, folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "folder not found", common.NotFoundError)
	}

	subfiles, err := dao.GetFilesByFolderId(ownerId, folder.ID)
	for _, subfile := range subfiles {
		if file.Name == subfile.Name {
			return errors.New(config.BottleSrvName, "Already a file with same name in parent folder", common.ConflictError)
		}
	}

	file.Name = name
	file.FolderId = folderId

	return dao.UpdateFiles(file)
}

func (f *FileHandler) RemoveFile(ctx context.Context, req *pb.RemoveFileRequest, resp *pb.RemoveFileResponse) error {
	file, err := dao.GetFileById(req.GetOwnerId(), req.GetFileId())
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	file.Status = db.Invalid
	now := time.Now()
	file.DeleteTime = &now

	return dao.UpdateFiles(file)
}
