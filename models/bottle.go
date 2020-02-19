package models

type UserBottle struct {
	BasicModel
	UserID   int     `gorm:"not null"`
	RootID   string  `gorm:"varchar(64);not null"`
	Capacity float64 `gorm:"type:double;not null;default:1024"` // mb
	Remain   float64 `gorm:"type:double;not null;default:1024"` // mb
}

type DirectoryInfo struct {
	BasicModel
	Name     string `gorm:"type:varchar(128);not null"`
	Path     string `gorm:"type:text;not null"`
	BottleID int    `gorm:"not null"`
	ParentID int    `gorm:"not null"`
	OwnerID  int    `gorm:"not null"`
}

type FileInfo struct {
	BasicModel
	Name     string  `gorm:"type:varchar(128);not null"`
	FileID   string  `gorm:"type:varchar(64);not null"`
	Size     float64 `gorm:"type:double;not null"`
	BottleID int     `gorm:"not null"`
	FolderID int     `gorm:"not null"`
	OwnerID  int     `gorm:"not null"`
}
