package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
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

func GetFolderById(ownerId int64, id int64) (*FolderInfo, error) {
	var folder FolderInfo
	if err := db.DB().Where("id = ? AND owner_id = ? AND status = ?", id, ownerId, db.Valid).First(&folder).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &folder, nil
}

func GetFoldersByPath(ownerId int64, names []string) (*FolderInfo, error) {
	if len(names) == 0 {
		return nil, nil
	}
	var folder FolderInfo
	parentId := RootId
	for _, name := range names {
		folder = FolderInfo{}
		if err := db.DB().Debug().Where("parent_id = ? AND owner_id = ? AND name = ? AND status = ?",
			parentId, ownerId, name, db.Valid).First(&folder).Error; err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return nil, nil
			}
			log.WithError(err).Errorf("DB error")
			return nil, common.NewDBError(err)
		}
		parentId = folder.ID
	}
	return &folder, nil
}

func GetFoldersByParentId(ownerId, parentId int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("parent_id = ? AND owner_id = ? AND status = ?", parentId, ownerId, db.Valid).Find(&folders).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByParentIds(ownerId int64, parentIds []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("parent_id in (?) AND owner_id = ? AND status = ?", parentIds, ownerId, db.Valid).Find(&folders).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByIds(ownerId int64, ids []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("id in (?) AND owner_id = ? AND status = ?", ids, ownerId, db.Valid).Find(&folders).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func CreateFolder(folder *FolderInfo) error {
	if err := db.DB().Create(folder).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}

func UpdateFolder(folder *FolderInfo) error {
	if err := db.DB().Save(folder).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}

func UpdateFolderAndSub(folder *FolderInfo, folders []*FolderInfo, files []*FileInfo) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		for _, f := range folders {
			if err := tx.Save(f).Error; err != nil {
				log.WithError(err).Errorf("DB error")
				return common.NewDBError(err)
			}
		}

		for _, f := range files {
			if err := tx.Save(f).Error; err != nil {
				log.WithError(err).Errorf("DB error")
				return common.NewDBError(err)
			}
		}

		if err := tx.Save(folder).Error; err != nil {
			log.WithError(err).Errorf("DB error")
			return common.NewDBError(err)
		}
		return nil
	})
}
