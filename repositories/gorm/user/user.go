package user

import (
	"coffee-helper/models"
	"coffee-helper/repositories/gorm/model"
	"gorm.io/gorm"
)

type User struct {
	*model.Model[models.User]

	totalStaffs  uint
	totalEditors uint
	totalAdmins  uint
}

func New(db *gorm.DB) (*User, error) {
	model, err := model.New[models.User](db)
	if err != nil {
		return nil, err
	}

	u := &User{
		Model: model,
	}

	return u, nil
}

func (u *User) ByTgID(tgID int64) (*models.User, error) {
	return u.ByColumn("tg_id", tgID)
}
