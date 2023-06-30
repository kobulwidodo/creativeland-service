package entity

import (
	"go-clean/src/lib/auth"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string
	Password string `json:"-"`
	Nama     string
	IsAdmin  bool
	UmkmID   uint
	GuestID  string `json:"-" gorm:"-:all"`
}

type CreateUserParam struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Nama     string `binding:"required"`
	UmkmID   uint   `binding:"required"`
}

type LoginUserParam struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

func (u *User) ConvertToAuthUser() auth.User {
	return auth.User{
		ID:       u.ID,
		Username: u.Username,
		Password: u.Password,
		IsAdmin:  u.IsAdmin,
		GuestID:  u.GuestID,
		UmkmID:   u.UmkmID,
	}
}
