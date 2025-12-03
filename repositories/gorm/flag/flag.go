package flag

import (
	"coffee-helper/models"
	"coffee-helper/repositories/gorm/model"
	"gorm.io/gorm"
)

type Flag struct {
	*model.Model[models.Flag]
}

func New(db *gorm.DB) (*Flag, error) {
	model, err := model.New[models.Flag](db)
	if err != nil {
		return nil, err
	}

	u := &Flag{
		Model: model,
	}

	return u, nil
}
