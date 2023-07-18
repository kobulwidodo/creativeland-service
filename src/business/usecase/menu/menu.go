package menu

import (
	"context"
	"errors"
	menuDom "go-clean/src/business/domain/menu"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(inputParam entity.CreateMenuParam, menuParam entity.MenuParam) (entity.Menu, error)
	GetAll(param entity.MenuParam) ([]entity.Menu, error)
	Get(params entity.MenuParam) (entity.Menu, error)
	Update(param entity.MenuParam, inputParam entity.UpdateMenuParam) error
	Delete(param entity.MenuParam) error
	ValidateMenu(ctx context.Context, menuID uint, user auth.UserAuthInfo) error
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

func (m *menu) Create(inputParam entity.CreateMenuParam, menuParam entity.MenuParam) (entity.Menu, error) {
	menu, err := m.menu.Create(entity.Menu{
		Name:        inputParam.Name,
		Description: inputParam.Description,
		Price:       inputParam.Price,
		UmkmID:      menuParam.UmkmID,
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

func (m *menu) Get(params entity.MenuParam) (entity.Menu, error) {
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

	if err := m.menu.Update(entity.MenuParam{
		ID: menu.ID,
	}, inputParam); err != nil {
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

func (m *menu) ValidateMenu(ctx context.Context, menuID uint, user auth.UserAuthInfo) error {
	if menuID == 0 {
		return errors.New("please provide menu id")
	}

	menu, err := m.menu.Get(entity.MenuParam{
		ID: menuID,
	})
	if err != nil {
		return err
	}

	if !user.User.IsAdmin && menu.UmkmID != user.User.UmkmID {
		return errors.New("unauthorized")
	}

	return nil
}
