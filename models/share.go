package models

import "time"

const (
	FILE = 1 + iota
	DIRECTORY
)

const (
	READ = 1 + iota
	WRITE
)

type ShareToken struct {
	BasicModel
	SharerID   int        `gorm:"not null"`
	Token      string     `gorm:"not null"`
	Action     int        `gorm:"not null;default:1"`
	Code       string     `gorm:"not null"`
	ExpireTime *time.Time `gorm:"not null"`
}

type ShareRef struct {
	BasicModel
	TokenID    int		  `gorm:"not null"`
	ObjectType int        `gorm:"not null"`
	ObjectID   int        `gorm:"not null"`
}