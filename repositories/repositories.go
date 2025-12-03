package repositories

import "coffee-helper/models"

type Repositories struct {
	User    User
	Flag    Flag
	Menu    Menu
	Product Product
}

type Model[T any] interface {
	Total() uint
	ByID(id uint) (*T, error)
	ByColumn(column string, value any) (*T, error)
	ByColumnList(column string, value any, offset int, limit int) ([]T, error)
	ByColumnContains(column string, substrs []string) (*T, error)
	ByColumnContainsList(column string, substrs []string, offset int, limit int) ([]T, error)
	CountByColumnContains(column string, substrs ...string) (int, error)
	List(offset int, limit int) ([]T, error)
	Create(*T) error
	Update(*T) error
	Delete(*T) error
}

type User interface {
	Model[models.User]
	ByTgID(tgID int64) (*models.User, error)
}

type Flag interface {
	Model[models.Flag]
}

type Menu interface {
	Model[models.Menu]
}

type Product interface {
	Model[models.Product]
}
