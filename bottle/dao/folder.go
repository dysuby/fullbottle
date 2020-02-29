package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"time"
)

type FolderInfo struct {
	db.BasicModel
	Name     string `gorm:"type:varchar(128);not null"`
	ParentId int64  `gorm:"not null"`
	OwnerId  int64  `gorm:"not null"`
}

func (FolderInfo) TableName() string {
	return "folder_info"
}

func VirtualRoot() *FolderInfo {
	folder := &FolderInfo{}
	folder.ID = RootId
	unix := time.Unix(0, 0)
	folder.CreateTime = &unix
	folder.UpdateTime = &unix
	return folder
}

func GetFolderById(ownerId int64, id int64) (*FolderInfo, error) {
	var folder FolderInfo
	if err := db.DB().Where("id = ? AND owner_id = ? AND status = ?", id, ownerId, db.Valid).First(&folder).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return &folder, nil
}

func GetFoldersByPath(ownerId int64, names []string, baseFolder int64, filterFolders []int64) (*FolderInfo, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var folder FolderInfo
	parentId := baseFolder
	for _, name := range names {
		folder = FolderInfo{}
		query := db.DB()
		if filterFolders != nil {
			query = query.Where("id in (?)", filterFolders)
			filterFolders = nil // only filter top level, used by share service
		}
		if err := db.DB().Where("parent_id = ? AND owner_id = ? AND name = ? AND status = ?",
			parentId, ownerId, name, db.Valid).First(&folder).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, nil
			}
			return nil, common.NewDBError(err)
		}
		parentId = folder.ID
	}
	return &folder, nil
}

func GetFoldersByParentId(ownerId, parentId int64, filterFolders []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	query := db.DB()
	if filterFolders != nil {
		query = query.Where("id in (?)", filterFolders)
	}
	if err := query.Where("parent_id = ? AND owner_id = ? AND status = ?", parentId, ownerId, db.Valid).Find(&folders).Error; err != nil {
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByParentIds(ownerId int64, parentIds []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("parent_id in (?) AND owner_id = ? AND status = ?", parentIds, ownerId, db.Valid).Find(&folders).Error; err != nil {
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByIds(ownerId int64, ids []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("id in (?) AND owner_id = ? AND status = ?", ids, ownerId, db.Valid).Find(&folders).Error; err != nil {
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func CreateFolder(folder *FolderInfo) error {
	if err := db.DB().Create(folder).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func UpdateFolder(folder *FolderInfo) error {
	if err := db.DB().Save(folder).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func RemoveFolderAndSub(folder *FolderInfo, folders []*FolderInfo, files []*FileInfo) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		size := int64(0)
		now := time.Now()

		for _, f := range folders {
			f.Status = db.Invalid
			f.DeleteTime = &now
			if err := tx.Save(f).Error; err != nil {
				return common.NewDBError(err)
			}
		}

		for _, f := range files {
			f.Status = db.Invalid
			f.DeleteTime = &now
			if err := tx.Save(f).Error; err != nil {
				return common.NewDBError(err)
			}
			size += f.Size
		}

		folder.Status = db.Invalid
		folder.DeleteTime = &now
		if err := tx.Save(folder).Error; err != nil {
			return common.NewDBError(err)
		}

		var bottle BottleMeta
		if err := tx.Where("user_id = ? AND status = ?", folder.OwnerId, db.Valid).First(&bottle).Error; err != nil {
			return common.NewDBError(err)
		}

		bottle.Remain += size
		if err := tx.Save(&bottle).Error; err != nil {
			return err
		}

		return nil
	})
}
