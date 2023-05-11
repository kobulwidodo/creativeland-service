package cart

import (
	"context"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	menuDom "go-clean/src/business/domain/menu"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(ctx context.Context, params entity.CreateCartParam) (entity.Cart, error)
	GetListByUser(ctx context.Context) ([]entity.Cart, error)
	Delete(ctx context.Context, param entity.CartParam) error
	ValidateCart(ctx context.Context, cartId uint, guestId string) error
}

type cart struct {
	cart cartDom.Interface
	auth auth.Interface
	menu menuDom.Interface
}

func Init(cd cartDom.Interface, auth auth.Interface, md menuDom.Interface) Interface {
	c := &cart{
		cart: cd,
		auth: auth,
		menu: md,
	}

	return c
}

func (c *cart) Create(ctx context.Context, params entity.CreateCartParam) (entity.Cart, error) {
	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return entity.Cart{}, err
	}

	menu, err := c.menu.Get(entity.MenuParam{
		ID:     params.MenuID,
		UmkmID: params.UmkmID,
	})
	if err != nil {
		return entity.Cart{}, err
	}

	cart, err := c.cart.Create(entity.Cart{
		UmkmID:       params.UmkmID,
		MenuID:       params.MenuID,
		GuestID:      user.User.GuestID,
		Amount:       params.Amount,
		Status:       entity.StatusActive,
		TotalPrice:   params.Amount * menu.Price,
		PricePerItem: menu.Price,
	})
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) GetListByUser(ctx context.Context) ([]entity.Cart, error) {
	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return []entity.Cart{}, err
	}

	cart, err := c.cart.GetList(entity.CartParam{
		GuestID: user.User.GuestID,
		Status:  entity.StatusActive,
	})
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) ValidateCart(ctx context.Context, cartId uint, guestId string) error {
	cart, err := c.cart.Get(entity.CartParam{
		ID: cartId,
	})
	if err != nil {
		return err
	}

	if cart.GuestID != guestId {
		return errors.New("unauthorized")
	}

	return nil
}

func (c *cart) Delete(ctx context.Context, param entity.CartParam) error {
	if err := c.cart.Delete(param); err != nil {
		return err
	}

	return nil
}
