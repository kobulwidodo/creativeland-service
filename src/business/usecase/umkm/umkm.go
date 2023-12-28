package umkm

import (
	"context"
	"errors"
	umkmDom "go-clean/src/business/domain/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(params entity.CreateUmkmParam) (entity.Umkm, error)
	GetList(param entity.UmkmParam) ([]entity.Umkm, error)
	Get(params entity.UmkmParam) (entity.Umkm, error)
	Update(param entity.UmkmParam, inputParam entity.UpdateUmkmParam) error
	Delete(param entity.UmkmParam) error
	ValidateUmkm(ctx context.Context, umkmId uint, user auth.UserAuthInfo) error
	SaveImage(ctx context.Context, param entity.UmkmParam, fileLocation string) error
}

type umkm struct {
	umkm umkmDom.Interface
}

func Init(ud umkmDom.Interface) Interface {
	u := &umkm{
		umkm: ud,
	}

	return u
}

func (u *umkm) Create(params entity.CreateUmkmParam) (entity.Umkm, error) {
	umkm, err := u.umkm.Create(entity.Umkm{
		Name:             params.Name,
		Slogan:           params.Slogan,
		Status:           entity.StatusClose,
		OwnerName:        params.OwnerName,
		OwnerPhoneNumber: params.OwnerPhoneNumber,
	})
	if err != nil {
		return umkm, err
	}

	return umkm, nil
}

func (u *umkm) GetList(param entity.UmkmParam) ([]entity.Umkm, error) {
	umkms, err := u.umkm.GetList(param)
	if err != nil {
		return umkms, err
	}

	return umkms, nil
}

func (u *umkm) Get(params entity.UmkmParam) (entity.Umkm, error) {
	umkm, err := u.umkm.Get(params)
	if err != nil {
		return umkm, err
	}

	return umkm, nil
}

func (u *umkm) Update(param entity.UmkmParam, inputParam entity.UpdateUmkmParam) error {
	umkm, err := u.umkm.Get(param)
	if err != nil {
		return err
	}

	if err := u.umkm.Update(entity.UmkmParam{ID: umkm.ID}, inputParam); err != nil {
		return err
	}

	return nil
}

func (u *umkm) Delete(param entity.UmkmParam) error {
	if err := u.umkm.Delete(param); err != nil {
		return err
	}

	return nil
}

func (u *umkm) ValidateUmkm(ctx context.Context, umkmId uint, user auth.UserAuthInfo) error {
	if umkmId == 0 {
		return errors.New("please provide umkm id")
	}

	if umkmId != user.User.UmkmID && !user.User.IsAdmin {
		return errors.New("unauthorized")
	}

	return nil
}

func (u *umkm) SaveImage(ctx context.Context, param entity.UmkmParam, fileLocation string) error {
	umkm, err := u.umkm.Get(entity.UmkmParam{
		ID: param.ID,
	})
	if err != nil {
		return err
	}

	if err := u.umkm.Update(entity.UmkmParam{ID: umkm.ID}, entity.UpdateUmkmParam{
		ImgPath: fileLocation,
	}); err != nil {
		return err
	}

	return nil
}
