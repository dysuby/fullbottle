package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"time"
)

const (
	File = 1 + iota
	Folder
)

const (
	Public = 1 + iota
	Private
)

const (
	View = 1 + iota
	Download
)

type ShareInfo struct {
	db.BasicModel
	SharerId       int64       `gorm:"not null"`
	Token          string      `gorm:"unique;not null"` // identifier
	Code           string      `gorm:"not null"`        // access code, empty means public share
	Privacy        int32       `gorm:"type:smallint;not null;default:1"`
	ParentFolderId int64       `gorm:"not null"` // all share objects' parent
	ExpireTime     *time.Time  `gorm:"not null"`
	ShareRefs      []*ShareRef `gorm:"foreignkey:ShareId;save_associations:false;preload:false"`
}

type ShareRef struct {
	db.BasicModel
	ShareId   int64 `gorm:"not null"`
	EntryType int   `gorm:"not null"`
	EntryId   int64 `gorm:"not null"`
}

type ShareMetrics struct {
	db.BasicModel
	ShareId  int64 `gorm:"not null"`
	ViewerId int64 `gorm:"not null"` // viewer id, who is access the share
	Action   int32 `gorm:"type:smallint;not null"`
}

type MetricsResult struct {
	Action int32  `gorm:"action"`
	Times  int64  `gorm:"times"`
}

func InitShare(info *ShareInfo) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(info).Error; err != nil {
			return common.NewDBError(err)
		}

		for _, ref := range info.ShareRefs {
			if err := tx.Create(ref).Error; err != nil {
				return common.NewDBError(err)
			}
		}

		return nil
	})
}

func UpdateShare(info *ShareInfo) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(info).Error; err != nil {
			return common.NewDBError(err)
		}

		for _, ref := range info.ShareRefs {
			if err := tx.Save(ref).Error; err != nil {
				return common.NewDBError(err)
			}
		}

		return nil
	})
}

func CancelShare(info *ShareInfo) error {
	info.Status = db.Canceled
	now := time.Now()
	info.DeleteTime = &now
	if err := db.DB().Save(info).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func GetShareById(sharerId int64, id int64) (*ShareInfo, error) {
	var info *ShareInfo
	if err := db.DB().Where("id = ? AND share_id = ? AND status = ?", id, sharerId, db.Valid).First(&sharerId).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return info, nil
}

func GetShareByToken(token string) (*ShareInfo, error) {
	var info *ShareInfo
	if err := db.DB().Where("token = ? AND status = ?", token, db.Valid).Find(info).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return info, nil
}

func CreateShareMetrics(metric *ShareMetrics) error {
	if err := db.DB().Create(metric).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func GetShareMetrics(shareId int64) ([]MetricsResult, error) {
	var res []MetricsResult
	err := db.DB().Select("action, count(distinct viewer_id) AS times").
		Where("share_id = ? AND status = ?", shareId, db.Valid).Group("action").Find(res).Error
	if err != nil {
		return res, common.NewDBError(err)
	}
	return res, nil
}

func GetShareEntry(shareId int64) ([]*ShareRef, error) {
	var entris []*ShareRef
	if err := db.DB().Where("share_id = ? AND status = ?", shareId, db.Valid).Find(entris).Error; err != nil {
		return nil, common.NewDBError(err)
	}
	return entris, nil
}
