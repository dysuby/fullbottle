package db

import (
	"time"
)

const (
	Invalid = 1 + iota
	Valid
	Expired

	// file upload status
	Uploading
	Failed
)

type Fields map[string]interface{}

type BasicModel struct {
	ID         int64      `gorm:"primary_key;auto_increment" json:"id"`
	Status     int32      `gorm:"type:smallint;default:2" json:"status"`
	CreateTime *time.Time `gorm:"column:create_time;default:current_timestamp;not null" json:"createTime"`
	UpdateTime *time.Time `gorm:"column:update_time;default:current_timestamp on update current_timestamp;not null" json:"updateTime"`
	DeleteTime *time.Time `json:"column:delete_time;deleteTime,omitempty"`
}
