package db

import (
	"time"
)

const (
	INVALID = 1 + iota
	VALID
	EXPIRED
)

type Fields map[string]interface{}

type BasicModel struct {
	ID         int64      `gorm:"primary_key;auto_increment"`
	Status     int32      `gorm:"type:smallint;default:2"`
	CreateTime *time.Time `gorm:"default:current_timestamp;not null"`
	UpdateTime *time.Time `gorm:"default:current_timestamp on update current_timestamp;not null"`
	DeleteTime *time.Time
}
