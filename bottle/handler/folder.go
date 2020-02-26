package handler

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/v2/errors"
	"github.com/vegchic/fullbottle/bottle/dao"
	pb "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/bottle/util"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/cache"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/config"
	"time"
)

const FolderLockKey = "lock:folder_id=%d"

type FolderHandler struct{}

func (*FolderHandler) GetFolderInfo(ctx context.Context, req *pb.GetFolderInfoRequest, resp *pb.GetFolderInfoResponse) error {
	ownerId := req.GetOwnerId()

	var folder *dao.FolderInfo
	var err error
	switch req.GetIdent().(type) {
	case *pb.GetFolderInfoRequest_FolderId:
		folder, err = dao.GetFolderById(ownerId, req.GetFolderId())
	case *pb.GetFolderInfoRequest_Path:
		names := util.SplitPath(req.GetPath())
		if len(names) == 0 {
			folder = &dao.FolderInfo{}
			folder.ID = dao.RootId
			unix := time.Unix(0, 0)
			folder.CreateTime = &unix
			folder.UpdateTime = &unix
			break
		}
		folder, err = dao.GetFoldersByPath(ownerId, names)
	default:

	}

	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	folders, files, err := util.GetSubEntry(ownerId, folder.ID)
	if err != nil {
		return err
	}

	f := &pb.FolderInfo{
		FolderId:   folder.ID,
		Name:       folder.Name,
		ParentId:   folder.ParentId,
		CreateTime: folder.CreateTime.Unix(),
		UpdateTime: folder.UpdateTime.Unix(),
	}

	f.Files = make([]*pb.FileInfo, len(files))
	for i, v := range files {
		f.Files[i] = &pb.FileInfo{
			Id:         v.ID,
			FileId:     v.FileId,
			Name:       v.Name,
			Size:       v.Metadata.Size,
			Hash:       v.Metadata.Hash,
			FolderId:   v.FolderId,
			OwnerId:    v.OwnerId,
			CreateTime: v.CreateTime.Unix(),
			UpdateTime: v.UpdateTime.Unix(),
		}
	}
	f.Folders = make([]*pb.FolderInfo, len(folders))
	for i, v := range folders {
		f.Folders[i] = &pb.FolderInfo{
			FolderId:   v.ID,
			Name:       v.Name,
			ParentId:   v.ParentId,
			CreateTime: v.CreateTime.Unix(),
			UpdateTime: v.UpdateTime.Unix(),
		}
	}

	resp.Folder = f
	return nil
}

func (*FolderHandler) CreateFolder(ctx context.Context, req *pb.CreateFolderRequest, resp *pb.CreateFolderResponse) error {
	name := req.GetName()
	parentId := req.GetParentId()
	ownerId := req.GetOwnerId()

	lock, err := cache.Obtain(fmt.Sprintf(FolderLockKey, parentId), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()

	if parentId != dao.RootId {
		parent, err := dao.GetFolderById(ownerId, parentId)
		if err != nil {
			return err
		} else if parent == nil {
			return errors.New(config.BottleSrvName, "Parent folder not found", common.NotFoundError)
		}
	}

	folders, err := dao.GetFoldersByParentId(ownerId, parentId)
	if err != nil {
		return err
	}

	for _, v := range folders {
		if v.Name == name {
			return errors.New(config.BottleSrvName, "There is a folder with same name in parent folder", common.ExistedError)
		}
	}

	folder := &dao.FolderInfo{
		Name:     name,
		ParentId: parentId,
		OwnerId:  ownerId,
	}
	err = dao.CreateFolder(folder)
	if err != nil {
		return err
	}

	resp.FolderId = folder.ID
	return nil
}

func (*FolderHandler) UpdateFolder(ctx context.Context, req *pb.UpdateFolderRequest, resp *pb.UpdateFolderResponse) error {
	folderId := req.GetFolderId()
	name := req.GetName()
	parentId := req.GetParentId()
	ownerId := req.GetOwnerId()

	lock, err := cache.Obtain(fmt.Sprintf(FolderLockKey, parentId), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()

	folder := &dao.FolderInfo{}
	ids := []int64{folderId}
	if parentId != dao.RootId {
		ids = append(ids, parentId)
	}
	fs, err := dao.GetFoldersByIds(ownerId, ids)
	if err != nil {
		return err
	} else if len(fs) != len(ids) {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	for _, f := range fs {
		if folderId == f.ID {
			folder = f
			break
		}
	}

	// check name
	subfolders, err := dao.GetFoldersByParentId(ownerId, parentId)
	if err != nil {
		return err
	}
	for _, v := range subfolders {
		if name == v.Name && folderId != v.ID {
			return errors.New(config.BottleSrvName, "There is a folder with same name in parent folder", common.ExistedError)
		}
	}

	// check parent_id
	folders, _, err := util.GetSubEntryRecursive(ownerId, folderId)
	for _, sub := range folders {
		if sub.ID == parentId {
			return errors.New(config.BottleSrvName, "Recursive structure", common.ConflictError)
		}
	}

	folder.Name = name
	folder.ParentId = parentId

	err = dao.UpdateFolder(folder)
	if err != nil {
		return err
	}

	return nil
}

func (*FolderHandler) RemoveFolder(ctx context.Context, req *pb.RemoveFolderRequest, resp *pb.RemoveFolderResponse) error {
	folderId := req.GetFolderId()
	ownerId := req.GetOwnerId()

	lock, err := cache.Obtain(fmt.Sprintf(FolderLockKey, folderId), 100*time.Millisecond)
	if err != nil {
		return err
	}
	defer lock.Release()

	folder, err := dao.GetFolderById(ownerId, folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	subfolders, subfiles, err := util.GetSubEntryRecursive(ownerId, folder.ID)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, v := range subfolders {
		v.Status = db.Invalid
		v.DeleteTime = &now
	}

	for _, v := range subfiles {
		v.Status = db.Invalid
		v.DeleteTime = &now
	}

	folder.Status = db.Invalid
	folder.DeleteTime = &now

	err = dao.UpdateFolderAndSub(folder, subfolders, subfiles)
	if err != nil {
		return err
	}

	return nil
}
