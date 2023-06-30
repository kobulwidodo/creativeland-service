package umkm

import (
	"fmt"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(umkm entity.Umkm) (entity.Umkm, error)
	GetList(param entity.UmkmParam) ([]entity.Umkm, error)
	GetListInByID(ids []uint) ([]entity.Umkm, error)
	Get(param entity.UmkmParam) (entity.Umkm, error)
	Update(selectParam entity.UmkmParam, updateParam entity.UpdateUmkmParam) error
	Delete(param entity.UmkmParam) error
}

type umkm struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	u := &umkm{
		db: db,
	}

	return u
}

func (u *umkm) Create(umkm entity.Umkm) (entity.Umkm, error) {
	if err := u.db.Create(&umkm).Error; err != nil {
		return umkm, err
	}

	return umkm, nil
}

func (u *umkm) GetList(param entity.UmkmParam) ([]entity.Umkm, error) {
	umkms := []entity.Umkm{}

	if err := u.db.Where(param).Where("name LIKE ?", fmt.Sprintf("%%%s%%", param.Name)).Find(&umkms).Error; err != nil {
		return umkms, err
	}

	return umkms, nil
}

func (u *umkm) GetListInByID(ids []uint) ([]entity.Umkm, error) {
	umkms := []entity.Umkm{}

	if err := u.db.Where(ids).Find(&umkms).Error; err != nil {
		return umkms, err
	}

	return umkms, nil
}

func (u *umkm) Get(param entity.UmkmParam) (entity.Umkm, error) {
	umkm := entity.Umkm{}

	if err := u.db.Where(param).First(&umkm).Error; err != nil {
		return umkm, err
	}

	return umkm, nil
}

func (u *umkm) Update(selectParam entity.UmkmParam, updateParam entity.UpdateUmkmParam) error {
	if err := u.db.Model(entity.Umkm{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (u *umkm) Delete(param entity.UmkmParam) error {
	if err := u.db.Where(param).Delete(&entity.Umkm{}).Error; err != nil {
		return err
	}

	return nil
}
