package user

import (
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(user entity.User) (entity.User, error)
	Get(param entity.UserParam) (entity.User, error)
}

type user struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	a := &user{
		db: db,
	}

	return a
}

func (a *user) Create(user entity.User) (entity.User, error) {
	if err := a.db.Create(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) Get(param entity.UserParam) (entity.User, error) {
	user := entity.User{}

	if err := a.db.Where(param).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}
