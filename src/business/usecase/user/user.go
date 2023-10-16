package user

import (
	"context"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	umkmDom "go-clean/src/business/domain/umkm"
	userDom "go-clean/src/business/domain/user"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"

	"golang.org/x/crypto/bcrypt"
)

type Interface interface {
	Create(params entity.CreateUserParam) (entity.User, error)
	Login(params entity.LoginUserParam) (string, error)
	GenerateGuestToken() (string, error)
	Get(param entity.UserParam) (entity.User, error)
	Me(ctx context.Context) (entity.User, error)
}

type user struct {
	user userDom.Interface
	auth auth.Interface
	cart cartDom.Interface
	umkm umkmDom.Interface
}

func Init(ad userDom.Interface, auth auth.Interface, cd cartDom.Interface, ud umkmDom.Interface) Interface {
	a := &user{
		user: ad,
		auth: auth,
		cart: cd,
		umkm: ud,
	}

	return a
}

func (a *user) Create(params entity.CreateUserParam) (entity.User, error) {
	user := entity.User{
		Username: params.Username,
		Nama:     params.Nama,
		UmkmID:   params.UmkmID,
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}

	user.Password = string(hashPass)

	newUser, err := a.user.Create(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}

func (a *user) Get(param entity.UserParam) (entity.User, error) {
	user, err := a.user.Get(entity.UserParam{
		ID: param.ID,
	})
	if err != nil {
		return user, err
	}

	return user, nil
}

func (a *user) Login(params entity.LoginUserParam) (string, error) {
	user, err := a.user.Get(entity.UserParam{
		Username: params.Username,
	})
	if err != nil {
		return "", err
	}

	if user.ID == 0 {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(params.Password)); err != nil {
		return "", errors.New("user tidak ditemukan atau password tidak sesuai")
	}

	token, err := a.auth.GenerateToken(user.ConvertToAuthUser())
	if err != nil {
		return "", err
	}

	return token, nil
}

func (a *user) GenerateGuestToken() (string, error) {
	token, err := a.auth.GenerateGuestToken()
	if err != nil {
		return "", err
	}

	return token, nil
}

func (u *user) Me(ctx context.Context) (entity.User, error) {
	user, err := u.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.User{}, err
	}

	me, err := u.user.Get(entity.UserParam{
		ID: user.User.ID,
	})
	if err != nil {
		return me, err
	}

	if me.UmkmID != 0 {
		umkm, err := u.umkm.Get(entity.UmkmParam{
			ID: me.UmkmID,
		})
		if err != nil {
			return me, err
		}
		me.UmkmStatus = umkm.Status
	}

	return me, nil
}
