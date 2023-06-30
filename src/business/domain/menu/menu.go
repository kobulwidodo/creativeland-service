package menu

import (
	"fmt"
	"go-clean/src/business/entity"

	"gorm.io/gorm"
)

type Interface interface {
	Create(menu entity.Menu) (entity.Menu, error)
	GetAll(param entity.MenuParam) ([]entity.Menu, error)
	GetListInByID(ids []int64) ([]entity.Menu, error)
	Get(param entity.MenuParam) (entity.Menu, error)
	Update(selectParam entity.MenuParam, updateParam entity.UpdateMenuParam) error
	Delete(param entity.MenuParam) error
}

type menu struct {
	db *gorm.DB
}

func Init(db *gorm.DB) Interface {
	m := &menu{
		db: db,
	}

	return m
}

func (m *menu) Create(menu entity.Menu) (entity.Menu, error) {
	if err := m.db.Create(&menu).Error; err != nil {
		return menu, err
	}

	return menu, nil
}

func (m *menu) GetAll(param entity.MenuParam) ([]entity.Menu, error) {
	menus := []entity.Menu{}

	if err := m.db.Where(param).Where("name LIKE ?", fmt.Sprintf("%%%s%%", param.Name)).Find(&menus).Error; err != nil {
		return menus, err
	}

	return menus, nil
}

func (m *menu) GetListInByID(ids []int64) ([]entity.Menu, error) {
	menus := []entity.Menu{}

	if err := m.db.Where(ids).Find(&menus).Error; err != nil {
		return menus, err
	}

	return menus, nil
}

func (m *menu) Get(param entity.MenuParam) (entity.Menu, error) {
	menu := entity.Menu{}

	if err := m.db.Where(param).First(&menu).Error; err != nil {
		return menu, err
	}

	return menu, nil
}

func (m *menu) Update(selectParam entity.MenuParam, updateParam entity.UpdateMenuParam) error {
	if err := m.db.Model(entity.Menu{}).Where(selectParam).Updates(updateParam).Error; err != nil {
		return err
	}

	return nil
}

func (m *menu) Delete(param entity.MenuParam) error {
	if err := m.db.Where(param).Delete(&entity.Menu{}).Error; err != nil {
		return err
	}

	return nil
}
