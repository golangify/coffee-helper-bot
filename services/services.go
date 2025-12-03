package services

import (
	"coffee-helper/config"
	"coffee-helper/repositories"
	"coffee-helper/services/user"
)

type Services struct {
	User    *user.Service
	Menu    any
	Product any
	Search  any
}

func New(config *config.Config, repositories *repositories.Repositories) *Services {
	s := &Services{
		User: user.New(config, repositories),
	}

	return s
}
