package menu

import (
	"coffee-helper/models"
	"coffee-helper/repositories/gorm/model"
	"gorm.io/gorm"
)

type Menu struct {
	*model.Model[models.Menu]
}

func New(db *gorm.DB) (*Menu, error) {
	model, err := model.New[models.Menu](db)
	if err != nil {
		return nil, err
	}

	u := &Menu{
		Model: model,
	}

	return u, nil
}
