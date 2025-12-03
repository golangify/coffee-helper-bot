package role

import (
	"errors"
	"fmt"

	"coffee-helper/models"
	"coffee-helper/repositories"

	"github.com/mazen160/go-random"
)

var ErrNotFound = errors.New("not found")

type RoleIssue struct {
	Secret    string
	Initiator *models.User
	RoleFlag  string
}

func newRoleIssue(initiator *models.User, roleFlag string) (*RoleIssue, error) {
	secret, err := random.String(16)
	if err != nil {
		return nil, err
	}
	return &RoleIssue{
		Secret:    secret,
		Initiator: initiator,
		RoleFlag:  roleFlag,
	}, nil
}

type Service struct {
	repositories *repositories.Repositories

	activeRoleIssues []*RoleIssue
}

func New(repositories *repositories.Repositories) *Service {
	s := &Service{
		repositories: repositories,
	}

	return s
}

func (s *Service) GetOrCreateIssue(initiator *models.User, roleFlag string) (*RoleIssue, error) {
	roleIssue, err := newRoleIssue(initiator, roleFlag)
	if err != nil {
		return nil, err
	}
	for _, activeRoleIssue := range s.activeRoleIssues {
		if activeRoleIssue.Initiator.ID == initiator.ID && activeRoleIssue.RoleFlag == roleFlag {
			return activeRoleIssue, nil
		}
	}
	s.activeRoleIssues = append(s.activeRoleIssues, roleIssue)
	return roleIssue, nil
}

func (s *Service) GrantIssue(user *models.User, secret string) (*RoleIssue, error) {
	for _, roleIssue := range s.activeRoleIssues {
		if roleIssue.Secret == secret {
			if !user.IsAdmin() && roleIssue.Initiator.ID == user.ID {
				return nil, errors.New("только администраторы могут выдавать роль самому себе")
			}
			oldUserFlags := user.Flags
			user.Flags.Set(roleIssue.RoleFlag)
			if err := s.repositories.User.Update(user); err != nil {
				user.Flags = oldUserFlags
				return roleIssue, err
			}
			return roleIssue, nil
		}
	}
	return nil, ErrNotFound
}

func (s *Service) MakeUserStaff(initiator *models.User, user *models.User) (*models.Flag, error) {
	return s.setUserFlag(initiator, user, models.FlagUserStaff)
}

func (s *Service) MakeUserEditor(initiator *models.User, user *models.User) (*models.Flag, error) {
	return s.setUserFlag(initiator, user, models.FlagUserEditor)
}

func (s *Service) MakeUserAdmin(initiator *models.User, user *models.User) (*models.Flag, error) {
	return s.setUserFlag(initiator, user, models.FlagUserAdmin)
}

func (s *Service) NewChange(initiator *models.User, user *models.User, sflag string, isAdd bool) (*models.Flag, error) {
	if isAdd {
		sflag = "+" + sflag
	} else {
		sflag = "-" + sflag
	}
	flag := &models.Flag{
		Flag:            sflag,
		TargetUserID:    user.ID,
		TargetUser:      user,
		InitiatorUserID: initiator.ID,
		InitiatorUser:   initiator,
	}

	if err := s.repositories.Flag.Create(flag); err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *Service) setUserFlag(initiator *models.User, user *models.User, sflag string) (*models.Flag, error) {
	if user.Flags.Has(sflag) {
		return nil, fmt.Errorf("пользователь уже имеет роль «%s»", models.UserFlagTitle[sflag])
	}

	oldUserFlags := user.Flags
	user.Flags.Set(sflag)
	if err := s.repositories.User.Update(user); err != nil {
		user.Flags = oldUserFlags
		return nil, err
	}

	flag, err := s.NewChange(initiator, user, sflag, true)
	if err != nil {
		return nil, err
	}

	return flag, nil
}

func (s *Service) removeUserFlag(initiator *models.User, user *models.User, sflag string) (*models.Flag, error) {
	if !user.Flags.Has(sflag) {
		return nil, fmt.Errorf("у пользователя нет роли «%s»", models.UserFlagTitle[sflag])
	}

	oldUserFlags := user.Flags
	user.Flags.Remove(sflag)
	if err := s.repositories.User.Update(user); err != nil {
		user.Flags = oldUserFlags
		return nil, err
	}

	flag, err := s.NewChange(initiator, user, sflag, false)
	if err != nil {
		return nil, err
	}

	return flag, nil
}
