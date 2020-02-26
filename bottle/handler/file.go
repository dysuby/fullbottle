package handler

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/bottle/util"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	"strings"
)

type FileHandler struct{}

func (f *FileHandler) GetFileInfo(ctx context.Context, req *pb.GetFileInfoRequest, resp *pb.GetFileInfoResponse) error {
	file, err := dao.GetFileById(req.GetFileId())

	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	resp.File = &pb.FileInfo{
		Id:         file.ID,
		FileId:     file.FileId,
		Name:       file.Name,
		Path:       file.Path,
		Level:      file.Level,
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

	file, err := dao.GetFileById(req.GetFileId())
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}
	folder, err := dao.GetFolderById(folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "folder not found", common.NotFoundError)
	}

	folders, files, err := util.GetSubEntry(folderId)
	if len(folders)+len(files) > config.FolderMaxSub {
		return errors.New(config.BottleSrvName, fmt.Sprintf("Cannnot create more than %s entry in a folder", config.FolderMaxSub), common.ExceedError)
	}

	file.Name = name
	file.FolderId = folderId
	file.Level = folder.Level + 1
	file.Path = strings.Join([]string{folder.Path, folder.Name}, "/") + "/"

	return dao.UpdateFiles(file)
}

func (f *FileHandler) RemoveFile(ctx context.Context, req *pb.RemoveFileRequest, resp *pb.RemoveFileResponse) error {
	file, err := dao.GetFileById(req.GetFileId())
	if err != nil {
		return err
	} else if file == nil {
		return errors.New(config.BottleSrvName, "File not found", common.NotFoundError)
	}

	file.Status = db.Invalid

	return dao.UpdateFiles(file)
}
