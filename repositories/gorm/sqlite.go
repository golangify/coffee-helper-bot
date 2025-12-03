package gorm

import (
	"coffee-helper/config"
	"coffee-helper/repositories"
	"coffee-helper/repositories/gorm/flag"
	"coffee-helper/repositories/gorm/menu"
	"coffee-helper/repositories/gorm/product"
	"coffee-helper/repositories/gorm/user"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func NewSQLite(cfg *config.Config) (*repositories.Repositories, error) {
	db, err := gorm.Open(sqlite.Open(cfg.Database), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// TODO: обернуть в цикл

	user, err := user.New(db)
	if err != nil {
		return nil, err
	}

	flag, err := flag.New(db)
	if err != nil {
		return nil, err
	}

	menu, err := menu.New(db)
	if err != nil {
		return nil, err
	}

	product, err := product.New(db)
	if err != nil {
		return nil, err
	}

	repo := &repositories.Repositories{
		User:    user,
		Flag:    flag,
		Menu:    menu,
		Product: product,
	}

	return repo, nil
}
