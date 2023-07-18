package cart

import (
	"context"
	"errors"
	cartDom "go-clean/src/business/domain/cart"
	menuDom "go-clean/src/business/domain/menu"
	umkmDom "go-clean/src/business/domain/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/lib/auth"
)

type Interface interface {
	Create(ctx context.Context, params entity.CreateCartParam) (entity.Cart, error)
	DecreaseItem(ctx context.Context, params entity.CartParam) error
	GetListByUser(ctx context.Context) ([]entity.Cart, error)
	Delete(ctx context.Context, param entity.CartParam) error
	ValidateCart(ctx context.Context, cartId uint, guestId string) error
	GetCartCount(ctx context.Context) (int, error)
}

type cart struct {
	cart cartDom.Interface
	auth auth.Interface
	menu menuDom.Interface
	umkm umkmDom.Interface
}

func Init(cd cartDom.Interface, auth auth.Interface, md menuDom.Interface, ud umkmDom.Interface) Interface {
	c := &cart{
		cart: cd,
		auth: auth,
		menu: md,
		umkm: ud,
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

	cartExist, _ := c.cart.Get(entity.CartParam{
		GuestID: user.User.GuestID,
		UmkmID:  params.UmkmID,
		MenuID:  params.MenuID,
		Status:  entity.StatusInCart,
	})

	if cartExist.ID != 0 {
		if err := c.cart.Update(entity.CartParam{
			GuestID: user.User.GuestID,
			UmkmID:  params.UmkmID,
			MenuID:  params.MenuID,
			Status:  entity.StatusInCart,
		}, entity.UpdateCartParam{
			Amount:     cartExist.Amount + params.Amount,
			TotalPrice: cartExist.TotalPrice + (menu.Price * params.Amount),
		}); err != nil {
			return cartExist, err
		}

		return cartExist, nil
	}

	cart, err := c.cart.Create(entity.Cart{
		UmkmID:       params.UmkmID,
		MenuID:       params.MenuID,
		GuestID:      user.User.GuestID,
		Amount:       params.Amount,
		Status:       entity.StatusInCart,
		TotalPrice:   params.Amount * menu.Price,
		PricePerItem: menu.Price,
	})
	if err != nil {
		return cart, err
	}

	return cart, nil
}

func (c *cart) DecreaseItem(ctx context.Context, params entity.CartParam) error {
	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return err
	}

	cart, err := c.cart.Get(entity.CartParam{
		ID:      params.ID,
		Status:  entity.StatusInCart,
		GuestID: user.User.GuestID,
	})
	if err != nil {
		return err
	}

	menu, err := c.menu.Get(entity.MenuParam{
		ID: cart.MenuID,
	})
	if err != nil {
		return err
	}

	if cart.Amount == 1 {
		if err := c.cart.Delete(entity.CartParam{
			ID:      params.ID,
			Status:  entity.StatusInCart,
			GuestID: user.User.GuestID,
		}); err != nil {
			return err
		}

		return nil
	}

	if err := c.cart.Update(entity.CartParam{
		ID:      params.ID,
		Status:  entity.StatusInCart,
		GuestID: user.User.GuestID,
	}, entity.UpdateCartParam{
		Amount:     cart.Amount - 1,
		TotalPrice: cart.TotalPrice - menu.Price,
	}); err != nil {
		return err
	}

	return nil
}

func (c *cart) GetListByUser(ctx context.Context) ([]entity.Cart, error) {
	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return []entity.Cart{}, err
	}

	cart, err := c.cart.GetList(entity.CartParam{
		GuestID: user.User.GuestID,
		Status:  entity.StatusInCart,
	})
	if err != nil {
		return cart, err
	}

	mapMenuIDs := make(map[uint]bool)
	mapUmkmIDs := make(map[uint]bool)
	for _, c := range cart {
		mapMenuIDs[c.MenuID] = true
		mapUmkmIDs[c.UmkmID] = true
	}

	menuIds := []int64{}
	for m := range mapMenuIDs {
		menuIds = append(menuIds, int64(m))
	}

	umkmIds := []uint{}
	for m := range mapUmkmIDs {
		umkmIds = append(umkmIds, m)
	}

	menus, err := c.menu.GetListInByID(menuIds)
	if err != nil {
		return cart, err
	}

	menusMap := make(map[uint]entity.Menu)
	for _, m := range menus {
		menusMap[m.ID] = m
	}

	umkms, err := c.umkm.GetListInByID(umkmIds)
	if err != nil {
		return cart, err
	}

	umkmsMap := make(map[uint]entity.Umkm)
	for _, u := range umkms {
		umkmsMap[u.ID] = u
	}

	for i, c := range cart {
		cart[i].Umkm = umkmsMap[c.UmkmID]
		cart[i].Menu = menusMap[c.MenuID]
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

func (c *cart) GetCartCount(ctx context.Context) (int, error) {
	result := 0

	user, err := c.auth.GetUserAuthInfo(ctx)
	if err != nil {
		return result, err
	}

	carts, err := c.cart.GetList(entity.CartParam{
		GuestID: user.User.GuestID,
		Status:  entity.StatusInCart,
	})
	if err != nil {
		return result, err
	}

	for _, c := range carts {
		result += c.Amount
	}

	return result, nil
}
