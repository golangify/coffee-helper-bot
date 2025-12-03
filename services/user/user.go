package user

import (
	"coffee-helper/config"
	"coffee-helper/models"
	"coffee-helper/repositories"
	"coffee-helper/services/pagination"
	"coffee-helper/services/user/role"
)

type UsersList struct {
	Users      []models.User
	Pagination *pagination.Pagination
}

type Service struct {
	config       *config.Config
	repositories *repositories.Repositories

	Role *role.Service
}

func New(config *config.Config, repositories *repositories.Repositories) *Service {
	s := &Service{
		config:       config,
		repositories: repositories,

		Role: role.New(repositories),
	}

	return s
}

func (s *Service) New(user *models.User) error {
	return s.repositories.User.Create(user)
}

func (s *Service) ByTgID(id int64) (*models.User, error) {
	return s.repositories.User.ByTgID(id)
}

func (s *Service) UsersCountByRolesFlags(flag ...string) (int, error) {
	count, err := s.repositories.User.CountByColumnContains("flags", flag...)
	if err != nil {
		return count, err
	}
	return count, nil
}

func (s *Service) AdminsList(page int) (*UsersList, error) {
	totalAdmins, err := s.UsersCountByRolesFlags(models.FlagUserAdmin)
	if err != nil {
		return nil, err
	}

	pagination, limit, offset := pagination.Paginate(page, s.config.ItemsPerPage, totalAdmins)

	admins, err := s.repositories.User.ByColumnContainsList("flags", []string{models.FlagUserAdmin}, offset, limit)
	if err != nil {
		return nil, err
	}

	usersList := &UsersList{
		Users:      admins,
		Pagination: pagination,
	}

	return usersList, nil
}
