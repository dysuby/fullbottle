package service

import (
	"context"
	pbbottle "github.com/vegchic/fullbottle/bottle/proto/bottle"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/share/dao"
)

func ValidateEntries(ctx context.Context, sharerId int64, parentId int64, folderIds []int64, fileIds []int64) ([]*dao.ShareRef, error) {
	var res []*dao.ShareRef

	for _, id := range folderIds {
		ref := &dao.ShareRef{ShareId: sharerId, EntryId: id, EntryType: dao.Folder}
		res = append(res, ref)
	}
	for _, id := range fileIds {
		ref := &dao.ShareRef{ShareId: sharerId, EntryId: id, EntryType: dao.File}
		res = append(res, ref)
	}

	bottleClient := common.BottleSrvClient()
	_, err := bottleClient.ValidateEntry(ctx, &pbbottle.ValidateEntryRequest{OwnerId: sharerId, ParentId: parentId, FileIds: fileIds, FolderIds: folderIds})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetShareFolder(ctx context.Context, info *dao.ShareInfo, path string) (*pbbottle.GetFolderInfoResponse, error) {
	entries, err := dao.GetShareEntry(info.ID)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return &pbbottle.GetFolderInfoResponse{Folder: &pbbottle.FolderInfo{}}, nil
	}

	var folderIds []int64
	var fileIds []int64
	for _, e := range entries {
		if e.EntryType == dao.Folder {
			folderIds = append(folderIds, e.EntryId)
		} else if e.EntryType == dao.File {
			fileIds = append(fileIds, e.EntryId)
		}
	}

	bottleClient := common.BottleSrvClient()

	// construct a request with base folder, filter files/folders
	folderReq := &pbbottle.GetFolderInfoRequest{OwnerId: info.SharerId,
		Ident: &pbbottle.GetFolderInfoRequest_Path_{Path: &pbbottle.GetFolderInfoRequest_Path{
			BaseFolder: info.ParentFolderId, Relative: path, FilterFiles: fileIds, FilterFolders: folderIds}}}

	folderResp, err := bottleClient.GetFolderInfo(ctx, folderReq)
	if err != nil {
		return nil, err
	}
	return folderResp, nil
}
