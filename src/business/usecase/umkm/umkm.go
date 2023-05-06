package umkm

import (
	umkmDom "go-clean/src/business/domain/umkm"
	"go-clean/src/business/entity"
)

type Interface interface {
	Create(params entity.CreateUmkmParam) (entity.Umkm, error)
	GetAll(param entity.UmkmParam) ([]entity.Umkm, error)
	GetById(params entity.UmkmParam) (entity.Umkm, error)
	Update(param entity.UmkmParam, inputParam entity.UpdateUmkmParam) error
	Delete(param entity.UmkmParam) error
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
		Name:   params.Name,
		Slogan: params.Slogan,
	})
	if err != nil {
		return umkm, err
	}

	return umkm, nil
}

func (u *umkm) GetAll(param entity.UmkmParam) ([]entity.Umkm, error) {
	umkms, err := u.umkm.GetAll(param)
	if err != nil {
		return umkms, err
	}

	return umkms, nil
}

func (u *umkm) GetById(params entity.UmkmParam) (entity.Umkm, error) {
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

	if err := u.umkm.Update(umkm, inputParam); err != nil {
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