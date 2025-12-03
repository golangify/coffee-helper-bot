package models

import (
	"gorm.io/gorm"
)

type Menu struct {
	gorm.Model
	UserID      uint
	User        *User  `gorm:"foreignKey:UserID"`
	Name        string `gorm:"notnull"`
	Description string
	ImageFileID string
}
