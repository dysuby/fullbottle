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
	Token          string      `gorm:"type:varchar(256);unique;not null"` // identifier
	Code           string      `gorm:"type:varchar(64);not null"`        // access code, empty means public share
	Privacy        int32       `gorm:"type:smallint;not null;default:1"`
	ParentFolderId int64       `gorm:"not null"` // all share objects' parent
	ExpireTime     *time.Time  `gorm:""`
	ShareRefs      []*ShareRef `gorm:"foreignkey:ShareId;save_associations:false;preload:false"`
}

func (ShareInfo) TableName() string {
	return "share_info"
}

type ShareRef struct {
	db.BasicModel
	ShareId   int64 `gorm:"not null"`
	EntryType int   `gorm:"type:smallint;not null"`
	EntryId   int64 `gorm:"not null"`
}

func (ShareRef) TableName() string {
	return "share_ref"
}

type ShareMetrics struct {
	db.BasicModel
	ShareId  int64 `gorm:"not null"`
	ViewerId int64 `gorm:"not null"` // viewer id, who access the share
	Action   int32 `gorm:"type:smallint;not null"`
}

func (ShareMetrics) TableName() string {
	return "share_metrics"
}

type MetricsResult struct {
	Action int32
	Times  int64
}

func InitShare(info *ShareInfo) error {
	return db.DB().Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(info).Error; err != nil {
			return common.NewDBError(err)
		}

		for _, ref := range info.ShareRefs {
			ref.ShareId = info.ID
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

func CancelShare(info *ShareInfo, status int32) error {
	if status != db.Canceled && status != db.Expired {
		return nil
	}
	info.Status = status
	now := time.Now()
	info.DeleteTime = &now
	if err := db.DB().Save(info).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func GetShareById(sharerId int64, id int64) (*ShareInfo, error) {
	var info ShareInfo
	if err := db.DB().Where("id = ? AND sharer_id = ? AND status = ?", id, sharerId, db.Valid).First(&info).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return &info, nil
}

func GetShareByToken(token string) (*ShareInfo, error) {
	var info ShareInfo
	if err := db.DB().Where("token = ? AND status = ?", token, db.Valid).Find(&info).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		return nil, common.NewDBError(err)
	}
	return &info, nil
}

func CreateShareMetrics(metric *ShareMetrics) error {
	if err := db.DB().Create(metric).Error; err != nil {
		return common.NewDBError(err)
	}
	return nil
}

func GetShareMetrics(shareId int64) ([]MetricsResult, error) {
	var res []MetricsResult
	err := db.DB().Table("share_metrics").Select("action, count(distinct viewer_id) AS times").
		Where("share_id = ? AND status = ?", shareId, db.Valid).Group("action").Scan(&res).Error
	if err != nil {
		return res, common.NewDBError(err)
	}
	return res, nil
}

func GetShareEntry(shareId int64) ([]*ShareRef, error) {
	var entries []*ShareRef
	if err := db.DB().Where("share_id = ? AND status = ?", shareId, db.Valid).Find(&entries).Error; err != nil {
		return nil, common.NewDBError(err)
	}
	return entries, nil
}
