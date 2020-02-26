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
	Path     string `gorm:"type:text;not null"`
	Level    int64  `gorm:"not null"`
	ParentID int64  `gorm:"not null"`
	OwnerId  int64  `gorm:"not null"`
}

func (FolderInfo) TableName() string {
	return "folder_info"
}

func GetFolderById(id int64) (*FolderInfo, error) {
	var folder FolderInfo
	if err := db.DB().Where("id = ? AND status = ?", id, db.Valid).First(&folder).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &folder, nil
}

func GetFolderOwner(id int64) (*FolderInfo, error) {
	var folder FolderInfo
	if err := db.DB().Select([]string{"owner_id"}).Where("id = ? AND status = ?", id, db.Valid).First(&folder).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &folder, nil
}

func GetFoldersByParentId(parentId int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("parent_id = ? AND status = ?", parentId, db.Valid).Find(&folders).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByParentIds(parentIds []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("parent_id in (?) AND status = ?", parentIds, db.Valid).Find(&folders).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return folders, nil
}

func GetFoldersByIds(ids []int64) ([]*FolderInfo, error) {
	var folders []*FolderInfo
	if err := db.DB().Where("id in (?) AND status = ?", ids, db.Valid).Find(&folders).Error; err != nil {
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
