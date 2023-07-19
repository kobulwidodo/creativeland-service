package cart_test

import (
	"context"
	mock_cart "go-clean/src/business/domain/mock/cart"
	mock_menu "go-clean/src/business/domain/mock/menu"
	mock_umkm "go-clean/src/business/domain/mock/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/cart"
	"go-clean/src/lib/auth"
	mock_auth "go-clean/src/lib/tests/mock/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_cart_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)

	createCartParamMock := entity.CreateCartParam{
		MenuID: 1,
		UmkmID: 1,
		Amount: 1,
	}

	userAuthMock := auth.UserAuthInfo{
		User: auth.User{
			GuestID: "1",
		},
	}

	menuParamMock := entity.MenuParam{
		ID:     1,
		UmkmID: 1,
	}

	menuResultMock := entity.Menu{
		Price: 10000,
	}

	cartParamMock := entity.CartParam{
		GuestID: "1",
		UmkmID:  1,
		MenuID:  1,
		Status:  entity.StatusInCart,
	}

	cartResultMock := entity.Cart{
		Model: gorm.Model{
			ID: 1,
		},
		Amount:     1,
		TotalPrice: 10000,
	}

	cartUpdateParamMock := entity.CartParam{
		GuestID: "1",
		UmkmID:  1,
		MenuID:  1,
		Status:  entity.StatusInCart,
	}

	cartUpdateMock := entity.UpdateCartParam{
		Amount:     2,
		TotalPrice: 20000,
	}

	createCartMock := entity.Cart{
		UmkmID:       1,
		MenuID:       1,
		GuestID:      "1",
		Amount:       1,
		Status:       entity.StatusInCart,
		TotalPrice:   10000,
		PricePerItem: 10000,
	}

	c := cart.Init(cartMock, authMock, menuMock, nil)

	type mockFields struct {
		auth *mock_auth.MockInterface
		menu *mock_menu.MockInterface
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
		menu: menuMock,
		cart: cartMock,
	}

	type args struct {
		ctx    context.Context
		params entity.CreateCartParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.Cart
		wantErr  bool
	}{
		{
			name: "failed to get user auth info",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to get menu",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(entity.Menu{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to update cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(assert.AnError)
			},
			want:    cartResultMock,
			wantErr: true,
		},
		{
			name: "all success to update cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().Update(cartUpdateParamMock, cartUpdateMock).Return(nil)
			},
			want:    cartResultMock,
			wantErr: false,
		},
		{
			name: "failed to create cart",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, nil)
				mock.cart.EXPECT().Create(createCartMock).Return(entity.Cart{}, assert.AnError)
			},
			want:    entity.Cart{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:    context.Background(),
				params: createCartParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, nil)
				mock.cart.EXPECT().Create(createCartMock).Return(createCartMock, nil)
			},
			want:    createCartMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := c.Create(tt.args.ctx, tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cart_DecreaseItem(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)

	paramMock := entity.CartParam{
		ID: 1,
	}

	authUserMock := auth.UserAuthInfo{
		User: auth.User{
			GuestID: "1",
		},
	}

	cartParamMock := entity.CartParam{
		ID:      1,
		Status:  entity.StatusInCart,
		GuestID: "1",
	}

	cartResultMock := entity.Cart{
		MenuID:     1,
		Amount:     1,
		TotalPrice: 10000,
	}

	cartTwoAmountResultMock := entity.Cart{
		MenuID:     1,
		Amount:     2,
		TotalPrice: 20000,
	}

	menuParamMock := entity.MenuParam{
		ID: 1,
	}

	menuResultMock := entity.Menu{
		Price: 10000,
	}

	updateCartParamMock := entity.UpdateCartParam{
		Amount:     1,
		TotalPrice: 10000,
	}

	c := cart.Init(cartMock, authMock, menuMock, nil)

	type mockFields struct {
		auth *mock_auth.MockInterface
		menu *mock_menu.MockInterface
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
		menu: menuMock,
		cart: cartMock,
	}

	type args struct {
		ctx   context.Context
		param entity.CartParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get auth user",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to get cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to get menu",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(entity.Menu{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to delete cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Delete(cartParamMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to delete cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Delete(cartParamMock).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "failed to update cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartTwoAmountResultMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Update(cartParamMock, updateCartParamMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().Get(cartParamMock).Return(cartTwoAmountResultMock, nil)
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.cart.EXPECT().Update(cartParamMock, updateCartParamMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := c.DecreaseItem(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.DecreaseItem() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cart_GetListByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)
	umkmMock := mock_umkm.NewMockInterface(ctrl)

	authUserMock := auth.UserAuthInfo{
		User: auth.User{
			GuestID: "1",
		},
	}

	cartParamMock := entity.CartParam{
		GuestID: "1",
		Status:  entity.StatusInCart,
	}

	cartResultMock := []entity.Cart{
		{
			MenuID: 1,
			UmkmID: 1,
		},
	}

	menusIDsMock := []int64{1}

	menuResultMock := []entity.Menu{
		{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}

	umkmIDsMock := []uint{1}

	umkmResultMock := []entity.Umkm{
		{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}

	resultMock := []entity.Cart{
		{
			Umkm: entity.Umkm{
				Model: gorm.Model{
					ID: 1,
				},
			},
			Menu: entity.Menu{
				Model: gorm.Model{
					ID: 1,
				},
			},
			MenuID: 1,
			UmkmID: 1,
		},
	}

	c := cart.Init(cartMock, authMock, menuMock, umkmMock)

	type mockFields struct {
		auth *mock_auth.MockInterface
		cart *mock_cart.MockInterface
		menu *mock_menu.MockInterface
		umkm *mock_umkm.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
		cart: cartMock,
		menu: menuMock,
		umkm: umkmMock,
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     []entity.Cart
		wantErr  bool
	}{
		{
			name: "failed to get auth user info",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			want:    []entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to get cart list",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    []entity.Cart{},
			wantErr: true,
		},
		{
			name: "failed to get menu list",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID(menusIDsMock).Return([]entity.Menu{}, assert.AnError)
			},
			want:    cartResultMock,
			wantErr: true,
		},
		{
			name: "failed to get umkm list",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID(menusIDsMock).Return(menuResultMock, nil)
				mock.umkm.EXPECT().GetListInByID(umkmIDsMock).Return([]entity.Umkm{}, assert.AnError)
			},
			want:    cartResultMock,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID(menusIDsMock).Return(menuResultMock, nil)
				mock.umkm.EXPECT().GetListInByID(umkmIDsMock).Return(umkmResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := c.GetListByUser(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.GetListByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cart_GetCartCount(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)

	authUserMock := auth.UserAuthInfo{
		User: auth.User{
			GuestID: "1",
		},
	}

	cartParamMock := entity.CartParam{
		GuestID: "1",
		Status:  entity.StatusInCart,
	}

	cartResultMock := []entity.Cart{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Amount: 1,
		},
	}

	resultMock := 1

	c := cart.Init(cartMock, authMock, nil, nil)

	type mockFields struct {
		auth *mock_auth.MockInterface
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		auth: authMock,
		cart: cartMock,
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     int
		wantErr  bool
	}{
		{
			name: "failed to get auth user info",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to get cart list",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(authUserMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := c.GetCartCount(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.GetCartCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_cart_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartMock := mock_cart.NewMockInterface(ctrl)

	paramMock := entity.CartParam{
		ID: 1,
	}

	c := cart.Init(cartMock, nil, nil, nil)

	type mockFields struct {
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		cart: cartMock,
	}

	type args struct {
		ctx   context.Context
		param entity.CartParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to delete cart",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().Delete(paramMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:   context.Background(),
				param: paramMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().Delete(paramMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := c.Delete(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("cart.GetCartCount() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_cart_ValidateCart(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartMock := mock_cart.NewMockInterface(ctrl)

	cartParamMock := entity.CartParam{
		ID: 1,
	}

	cartResultMock := entity.Cart{
		GuestID: "1",
	}

	c := cart.Init(cartMock, nil, nil, nil)

	type mockFields struct {
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		cart: cartMock,
	}

	type args struct {
		ctx     context.Context
		cartID  uint
		guestID string
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get cart",
			args: args{
				ctx:     context.Background(),
				cartID:  1,
				guestID: "1",
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().Get(cartParamMock).Return(entity.Cart{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed unauthorized",
			args: args{
				ctx:     context.Background(),
				cartID:  1,
				guestID: "2",
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:     context.Background(),
				cartID:  1,
				guestID: "1",
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().Get(cartParamMock).Return(cartResultMock, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := c.ValidateCart(tt.args.ctx, tt.args.cartID, tt.args.guestID)
			if (err != nil) != tt.wantErr {
				t.Errorf("menu.ValidateMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
