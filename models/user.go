package models

import (
	"coffee-helper/models/flags"
	"gorm.io/gorm"
	"strings"
)

const (
	FlagUserAdmin  = "a"
	FlagUserEditor = "e"
	FlagUserStaff  = "b"
	FlagUserBanned = "B"
)

var UserFlagTitle = map[string]string{
	FlagUserAdmin:  "администратор",
	FlagUserEditor: "редактор",
	FlagUserStaff:  "сотрудник",
	FlagUserBanned: "забаненный",
}

type User struct {
	gorm.Model
	TgID      int64 `gorm:"unique;notnull"`
	Flags     flags.Flags
	FirstName string
	LastName  string
	Username  string
}

func (u *User) IsAdmin() bool {
	return u.Flags.Has(FlagUserAdmin)
}

func (u *User) IsEditor() bool {
	return u.IsAdmin() || u.Flags.Has(FlagUserEditor)
}

func (u *User) IsStaff() bool {
	return u.IsAdmin() || u.IsEditor() || u.Flags.Has(FlagUserStaff)
}

func (u *User) IsBanned() bool {
	return u.Flags.Has(FlagUserBanned)
}

func (u *User) RolesRu() []string {
	var res []string
	for _, flag := range u.Flags {
		res = append(res, UserFlagTitle[string(flag)])
	}
	return res
}

func (u *User) String() string {
	var b strings.Builder
	b.WriteString(u.FirstName)
	if u.LastName != "" {
		b.WriteString(" " + u.LastName)
	}
	if u.Username != "" {
		b.WriteString(" (@" + u.Username + ")")
	}
	return b.String()
}
