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

// return id=0 when file not found, this rpc just for upload check
func (f *FileHandler) GetFileByMeta(ctx context.Context, req *pb.GetFileByMetaRequest, resp *pb.GetFileByMetaResponse) error {
	file, err := dao.GetFileByUploadMeta(req.GetOwnerId(), req.GetName(), req.GetFolderId(), req.GetMetaId())
	if err != nil {
		return err
	} else if file == nil {
		file = &dao.FileInfo{}
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

func (f *FileHandler) CreateFile(ctx context.Context, req *pb.CreateFileRequest, resp *pb.CreateFileResponse) error {
	meta, err := dao.GetFileMetaById(req.GetMetaId())
	if err != nil {
		return err
	} else if meta == nil {
		return errors.New(config.BottleSrvName, "Meta not found", common.NotFoundError)
	}

	info := &dao.FileInfo{
		FolderId: req.GetFolderId(),
		OwnerId:  req.GetOwnerId(),
		FileId:   req.GetMetaId(),
		Name:     req.GetName(),
	}

	err = service.CreateFile(info, meta)
	if err != nil {
		return err
	}

	resp.Id = info.ID
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

	subfiles, err := dao.GetFilesByFolderId(ownerId, folder.ID, nil)
	for _, subfile := range subfiles {
		if file.Name == subfile.Name && file.ID != subfile.ID {
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

	return dao.RemoveFile(file.OwnerId, file)
}

// return id=0 when meta not found, this rpc just for upload check
func (f *FileHandler) GetFileMeta(ctx context.Context, req *pb.GetFileMetaRequest, resp *pb.GetFileMetaResponse) error {
	hash := req.GetHash()

	meta, err := dao.GetFileMetaByHash(hash)
	if err != nil {
		return err
	} else if meta == nil {
		meta = &dao.FileMeta{}
	}

	resp.Id = meta.ID
	resp.Fid = meta.Fid
	resp.Hash = meta.Hash
	resp.Size = meta.Size
	resp.ChunkManifest = meta.ChunkManifest

	return nil
}

func (f *FileHandler) CreateFileMeta(ctx context.Context, req *pb.CreateFileMetaRequest, resp *pb.CreateFileMetaResponse) error {
	meta := &dao.FileMeta{
		Fid:           req.GetFid(),
		Hash:          req.GetHash(),
		Size:          req.GetSize(),
		ChunkManifest: req.GetChunkManifest(),
	}

	err := dao.CreateFileMeta(meta)
	if err != nil {
		return err
	}

	resp.Id = meta.ID
	return nil
}
