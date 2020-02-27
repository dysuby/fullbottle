package dao

import (
	"github.com/jinzhu/gorm"
	"github.com/vegchic/fullbottle/common"
	"github.com/vegchic/fullbottle/common/db"
	"github.com/vegchic/fullbottle/common/log"
)

const (
	B = 1 << (10 * iota)
	KB
	MB
	GB
)

const RootId int64 = -1

// user store space metadata
type BottleMeta struct {
	db.BasicModel
	UserID   int64 `gorm:"not null"`
	Capacity int64 `gorm:"type:bigint;not null"` // b
	Remain   int64 `gorm:"type:bigint;not null"` // b
}

func (BottleMeta) TableName() string {
	return "bottle_meta"
}

func GetBottlesById(id int64) (*BottleMeta, error) {
	var bottle BottleMeta
	if err := db.DB().Where("id = ? AND status = ?", id, db.Valid).First(&bottle).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &bottle, nil
}

func GetBottlesByUserId(uid int64) (*BottleMeta, error) {
	var bottle BottleMeta
	if err := db.DB().Where("user_id = ? AND status = ?", uid, db.Valid).First(&bottle).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return nil, nil
		}
		log.WithError(err).Errorf("DB error")
		return nil, common.NewDBError(err)
	}
	return &bottle, nil
}

func UpdateBottle(meta *BottleMeta) error {
	if err := db.DB().Save(meta).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}

func InitBottle(bottle *BottleMeta) error {
	if err := db.DB().Create(bottle).Error; err != nil {
		log.WithError(err).Errorf("DB error")
		return common.NewDBError(err)
	}
	return nil
}
