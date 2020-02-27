package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/weed"
)

type FileInfo struct {
	db.BasicModel
	Name     string   `gorm:"type:varchar(128);not null"`
	FileId   int64    `grom:"not null"`
	OwnerId  int64    `gorm:"not null"`
	FolderId int64    `gorm:"not null"`
	Metadata FileMeta `gorm:"foreignkey:FileId;save_associations:false;preload:false"`
}

func (FileInfo) TableName() string {
	return "file_info"
}

func (f *FileInfo) FromUploadMeta(meta *weed.FileUploadMeta) {
	f.Name = meta.Name
	f.OwnerId = meta.OwnerId
	f.FolderId = meta.FolderId
}

func GetFileById(ownerId, id int64) (*FileInfo, error) {
	var file FileInfo
	if err := db.DB().Where("id = ? AND owner_id = ? AND status = ?", id, ownerId, db.Valid).
		Preload("Metadata").First(&file).Error; err != nil {

		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &file, nil
}

func GetFilesByFolderId(ownerId, folderId int64) ([]*FileInfo, error) {
	// TODO check result
	var files []*FileInfo
	if err := db.DB().Where("folder_id = ? AND owner_id = ? AND status = ?", folderId, ownerId, db.Valid).
		Preload("Metadata").Find(&files).Error; err != nil {

		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return files, nil
}

func GetFilesByFolderIds(ownerId int64, parentIds []int64) ([]*FileInfo, error) {
	var files []*FileInfo
	if err := db.DB().Where("folder_id in (?) AND owner_id = ? AND status = ?", parentIds, ownerId, db.Valid).
		Preload("Metadata").Find(&files).Error; err != nil {

		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return files, nil
}

func UpdateFiles(file *FileInfo) error {
	if err := db.DB().Updates(file).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}

func CreateFile(file *FileInfo, meta *FileMeta) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(file).Error; err != nil {
			log.WithError(err).Errorf("DB error")
			return common.NewDBError(err)
		}

		var bottle BottleMeta
		if err := db.DB().Where("user_id = ? AND status = ?", file.OwnerId, db.Valid).First(&bottle).Error; err != nil {
			log.WithError(err).Errorf("DB error")
			return common.NewDBError(err)
		}

		bottle.Remain -= meta.Size
		if err := tx.Save(&bottle).Error; err != nil {
			return err
		}

		return nil
	})
}
