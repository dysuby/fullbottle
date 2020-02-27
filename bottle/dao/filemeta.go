package dao

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
	"github.com/vegchic/fullbottle/weed"
)

type FileMeta struct {
	db.BasicModel
	Fid  string `gorm:"type:varchar(64);not null"`
	Size int64  `gorm:"type:bigint;not null"`
	Hash string `gorm:"type:varchar(128);not null"`
	ChunkManifest string `gorm:"type:text;not null"`
}

func (FileMeta) TableName() string {
	return "file_meta"
}

func (f *FileMeta) FromUploadMeta(meta *weed.FileUploadMeta) {
	f.Fid = meta.Fid
	f.Size = meta.Size
	f.Hash = meta.Hash
	b, _ := meta.ChunkManifest.Json()
	f.ChunkManifest = string(b)
}

func (f *FileMeta) GetChunkInfo() (cm *weed.ChunkManifest, err error) {
	if err = json.Unmarshal([]byte(f.ChunkManifest), cm); err != nil {
		return
	}
	return
}

func GetFileMetaByHash(hash string) (*FileMeta, error) {
	var meta FileMeta
	if err := db.DB().Where("hash = ? AND status = ?", hash, db.Valid).First(&meta).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &meta, nil
}

func GetFileMetaById(id int64) (*FileMeta, error) {
	var meta FileMeta
	if err := db.DB().Where("id = ? AND status = ?", id, db.Valid).First(&meta).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &meta, nil
}

func CreateFileMeta(meta *FileMeta) error {
	if err := db.DB().Create(meta).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}
