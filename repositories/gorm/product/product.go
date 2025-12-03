package product

import (
	"coffee-helper/models"
	"coffee-helper/repositories/gorm/model"
	"gorm.io/gorm"
)

type Product struct {
	*model.Model[models.Product]
}

func New(db *gorm.DB) (*Product, error) {
	model, err := model.New[models.Product](db)
	if err != nil {
		return nil, err
	}

	u := &Product{
		Model: model,
	}

	return u, nil
}
