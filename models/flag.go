package models

import "gorm.io/gorm"

type Flag struct {
	gorm.Model
	Flag            string `gorm:"notnull"`
	TargetUserID    uint   `gorm:"notnull"`
	TargetUser      *User  `gorm:"foreignKey:TargetUserID"`
	InitiatorUserID uint   `gorm:"notnull"`
	InitiatorUser   *User  `gorm:"foreignKey:InitiatorUserID"`
}
