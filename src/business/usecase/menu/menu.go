package menu

import (
	menuDom "go-clean/src/business/domain/menu"
	"go-clean/src/business/entity"
)

type Interface interface {
	Create(params entity.CreateMenuParam) (entity.Menu, error)
	GetAll(param entity.MenuParam) ([]entity.Menu, error)
	GetById(params entity.MenuParam) (entity.Menu, error)
	Update(param entity.MenuParam, inputParam entity.UpdateMenuParam) error
	Delete(param entity.MenuParam) error
}

type menu struct {
	menu menuDom.Interface
}

func Init(md menuDom.Interface) Interface {
	m := &menu{
		menu: md,
	}

	return m
}

func (m *menu) Create(params entity.CreateMenuParam) (entity.Menu, error) {
	menu, err := m.menu.Create(entity.Menu{
		Name:        params.Name,
		Description: params.Description,
		Price:       params.Price,
		UmkmID:      params.UmkmID,
	})
	if err != nil {
		return menu, err
	}

	return menu, nil
}

func (m *menu) GetAll(param entity.MenuParam) ([]entity.Menu, error) {
	menus, err := m.menu.GetAll(param)
	if err != nil {
		return menus, err
	}

	return menus, nil
}

func (m *menu) GetById(params entity.MenuParam) (entity.Menu, error) {
	menu, err := m.menu.Get(params)
	if err != nil {
		return menu, err
	}

	return menu, nil
}

func (m *menu) Update(param entity.MenuParam, inputParam entity.UpdateMenuParam) error {
	menu, err := m.menu.Get(param)
	if err != nil {
		return err
	}

	if err := m.menu.Update(menu, inputParam); err != nil {
		return err
	}

	return nil
}

func (m *menu) Delete(param entity.MenuParam) error {
	if err := m.menu.Delete(param); err != nil {
		return err
	}

	return nil
}
