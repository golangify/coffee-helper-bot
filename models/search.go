package models

import (
	"gorm.io/gorm"
)

type Search struct {
	gorm.Model
	UserID uint
	User   *User `gorm:"foreignKey:UserID"`
	Text   string
}
