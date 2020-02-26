// TODO add lock for sub entry
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

type FolderHandler struct{}

func (FolderHandler) GetFolderInfo(ctx context.Context, req *pb.GetFolderInfoRequest, resp *pb.GetFolderInfoResponse) error {
	folderId := req.GetFolderId()

	folder, err := dao.GetFolderById(folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	folders, files, err := util.GetSubEntry(folderId)
	if err != nil {
		return err
	}

	f := pb.FolderInfo{
		FolderId:   folder.ID,
		Name:       folder.Name,
		Path:       folder.Path,
		ParentId:   folder.ParentID,
		CreateTime: folder.CreateTime.Unix(),
		UpdateTime: folder.UpdateTime.Unix(),
	}

	f.Files = make([]*pb.FileInfo, len(files))
	for i, v := range files {
		f.Files[i] = &pb.FileInfo{
			Id:         v.ID,
			FileId:     v.FileId,
			Name:       v.Name,
			Path:       v.Path,
			Size:       v.Metadata.Size,
			Hash:       v.Metadata.Hash,
			Level:      v.Level,
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
			Path:       v.Path,
			ParentId:   v.ParentID,
			Level:      v.Level,
			CreateTime: v.CreateTime.Unix(),
			UpdateTime: v.UpdateTime.Unix(),
		}
	}

	resp.Folder = &f
	return nil
}

func (FolderHandler) CreateFolder(ctx context.Context, req *pb.CreateFolderRequest, resp *pb.CreateFolderResponse) error {
	name := req.GetName()
	parentId := req.GetParentId()
	ownerId := req.GetOwnerId()

	parent, err := dao.GetFolderById(parentId)
	if err != nil {
		return err
	} else if parent == nil {
		return errors.New(config.BottleSrvName, "Parent folder not found", common.NotFoundError)
	}
	// simple check
	if parent.OwnerId != ownerId {
		return errors.New(config.BottleSrvName, "Owner cannot different from parent folder's", common.ConflictError)
	}

	if parent.Level+1 > config.FolderMaxLevel {
		return errors.New(config.BottleSrvName, fmt.Sprintf("Folder level cannot exceed %s", config.FolderMaxLevel), common.ExceedError)
	}

	path := strings.Join([]string{parent.Path, parent.Name}, "") + "/"
	level := parent.Level + 1

	folders, files, err := util.GetSubEntry(parentId)
	if err != nil {
		return err
	}
	for _, v := range folders {
		if v.Name == name {
			return errors.New(config.BottleSrvName, "There is a folder with same name in parent folder", common.ExistedError)
		}
	}

	if len(files)+len(folders) > config.FolderMaxSub {
		return errors.New(config.BottleSrvName, fmt.Sprintf("Cannnot create more than %s entry in a folder", config.FolderMaxSub), common.ExceedError)
	}

	folder := dao.FolderInfo{
		Name:     name,
		Path:     path,
		Level:    level,
		ParentID: parentId,
		OwnerId:  ownerId,
	}
	err = dao.CreateFolder(&folder)
	if err != nil {
		return err
	}

	resp.FolderId = folder.ID
	return nil
}

func (FolderHandler) UpdateFolder(ctx context.Context, req *pb.UpdateFolderRequest, resp *pb.UpdateFolderResponse) error {
	folderId := req.GetFolderId()
	name := req.GetName()
	parentId := req.GetParentId()

	var f, p *dao.FolderInfo
	fs, err := dao.GetFoldersByIds([]int64{folderId, parentId})
	if err != nil {
		return err
	} else if len(fs) != 2 {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	if fs[0].ID == folderId {
		f, p = fs[0], fs[1]
	} else {
		f, p = fs[1], fs[0]
	}

	if p.Level+1 > config.FolderMaxLevel {
		return errors.New(config.BottleSrvName, fmt.Sprintf("Folder level cannot exceed %s", config.FolderMaxLevel), common.ExceedError)
	}

	// 禁止套娃 & cannot remove root
	if strings.HasPrefix(p.Path, f.Path) {
		return errors.New(config.BottleSrvName, "Folder structure error", common.ConflictError)
	}

	newPath := strings.Join([]string{p.Path, p.Name}, "") + "/"
	newLevel := p.Level + 1

	subfolders, subfiles, err := util.GetSubEntry(parentId)
	if err != nil {
		return err
	}

	if len(subfiles)+len(subfolders) > config.FolderMaxSub {
		return errors.New(config.BottleSrvName, fmt.Sprintf("Cannnot create more than %s entry in a folder", config.FolderMaxSub), common.ExceedError)
	}

	for _, v := range subfolders {
		if name == v.Name {
			return errors.New(config.BottleSrvName, "There is a folder with same name in parent folder", common.ExistedError)
		}
	}

	subfolders, subfiles, err = util.GetSubEntryRecursive(f.ID)
	if err != nil {
		return err
	}

	newSubPath := strings.Join([]string{newPath, name}, "") + "/"
	oldSubPath := strings.Join([]string{f.Path, f.Name}, "") + "/"
	for _, v := range subfolders {
		v.Path = strings.Replace(v.Path, oldSubPath, newSubPath, 1)
		v.Level = newLevel + 1
	}

	for _, v := range subfiles {
		v.Path = strings.Replace(v.Path, oldSubPath, newSubPath, 1)
		v.Level = newLevel + 1
	}

	f.Path = newPath
	f.Name = name
	f.ParentID = parentId
	f.Level = newLevel

	err = dao.UpdateFolderAndSub(f, subfolders, subfiles)
	if err != nil {
		return err
	}

	return nil
}

func (FolderHandler) RemoveFolder(ctx context.Context, req *pb.RemoveFolderRequest, resp *pb.RemoveFolderResponse) error {
	folderId := req.GetFolderId()

	folder, err := dao.GetFolderById(folderId)
	if err != nil {
		return err
	} else if folder == nil {
		return errors.New(config.BottleSrvName, "Folder not found", common.NotFoundError)
	}

	if folder.Level == 0 {
		return errors.New(config.BottleSrvName, "Cannot remove root folder", common.ConflictError)
	}

	subfolders, subfiles, err := util.GetSubEntryRecursive(folder.ID)
	if err != nil {
		return err
	}

	for _, v := range subfolders {
		v.Status = db.Invalid
	}

	for _, v := range subfiles {
		v.Status = db.Invalid
	}

	folder.Status = db.Invalid

	err = dao.UpdateFolderAndSub(folder, subfolders, subfiles)
	if err != nil {
		return err
	}

	return nil
}
