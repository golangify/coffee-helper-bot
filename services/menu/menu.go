package user

import (
	"coffee-helper/config"
	"coffee-helper/models"
	"coffee-helper/repositories"
	"coffee-helper/services/pagination"
)

type MenusList struct {
	Menus      []models.Menu
	Pagination *pagination.Pagination
}

type Service struct {
	config       *config.Config
	repositories *repositories.Repositories
}

func New(config *config.Config, repositories *repositories.Repositories) *Service {
	s := &Service{
		config:       config,
		repositories: repositories,
	}
	return s
}

func (s *Service) New(creator *models.User, menu *models.Menu) error {
	menu.UserID = creator.ID
	menu.User = creator
	return s.repositories.Menu.Create(menu)
}
