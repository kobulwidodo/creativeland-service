package menu_test

import (
	"context"
	mock_menu "go-clean/src/business/domain/mock/menu"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/menu"
	"go-clean/src/lib/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_menu_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	createMenuParamMock := entity.CreateMenuParam{
		Name:        "menu",
		Description: "menu",
		Price:       10000,
	}

	menuParamMock := entity.MenuParam{
		UmkmID: 1,
	}

	newMenuMock := entity.Menu{
		Name:        createMenuParamMock.Name,
		Description: createMenuParamMock.Description,
		Price:       createMenuParamMock.Price,
		UmkmID:      menuParamMock.UmkmID,
	}

	menuResultMock := entity.Menu{
		Model: gorm.Model{
			ID: 1,
		},
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		inputParam entity.CreateMenuParam
		menuParam  entity.MenuParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.Menu
		wantErr  bool
	}{
		{
			name: "failed to create new menu",
			args: args{
				inputParam: createMenuParamMock,
				menuParam:  menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Create(newMenuMock).Return(entity.Menu{}, assert.AnError)
			},
			want:    entity.Menu{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				inputParam: createMenuParamMock,
				menuParam:  menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Create(newMenuMock).Return(menuResultMock, nil)
			},
			want:    menuResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(mocks, tt.args)
		got, err := m.Create(tt.args.inputParam, tt.args.menuParam)
		if (err != nil) != tt.wantErr {
			t.Errorf("menu.Create() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.want, got)
	}
}

func Test_menu_GetAll(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	menuParamMock := entity.MenuParam{}

	menuResultMock := []entity.Menu{
		{
			Name: "menu",
		},
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		param entity.MenuParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     []entity.Menu
		wantErr  bool
	}{
		{
			name: "failed to get all menu",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().GetAll(menuParamMock).Return([]entity.Menu{}, assert.AnError)
			},
			want:    []entity.Menu{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().GetAll(menuParamMock).Return(menuResultMock, nil)
			},
			want:    menuResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(mocks, tt.args)
		got, err := m.GetAll(tt.args.param)
		if (err != nil) != tt.wantErr {
			t.Errorf("menu.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.want, got)
	}
}

func Test_menu_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	menuParamMock := entity.MenuParam{
		ID: 1,
	}

	menuResultMock := entity.Menu{
		Name: "menu",
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		param entity.MenuParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.Menu
		wantErr  bool
	}{
		{
			name: "failed to get all menu",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(entity.Menu{}, assert.AnError)
			},
			want:    entity.Menu{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
			},
			want:    menuResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(mocks, tt.args)
		got, err := m.Get(tt.args.param)
		if (err != nil) != tt.wantErr {
			t.Errorf("menu.Get() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
		assert.Equal(t, tt.want, got)
	}
}

func Test_menu_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	menuParamMock := entity.MenuParam{
		ID: 1,
	}

	menuResultMock := entity.Menu{
		Model: gorm.Model{
			ID: 1,
		},
		Name: "menu",
	}

	updateMenuParamMock := entity.UpdateMenuParam{
		Name: "new menu",
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		param      entity.MenuParam
		inputParam entity.UpdateMenuParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get menu",
			args: args{
				param:      menuParamMock,
				inputParam: updateMenuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(entity.Menu{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update menu",
			args: args{
				param:      menuParamMock,
				inputParam: updateMenuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.menu.EXPECT().Update(menuParamMock, updateMenuParamMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				param:      menuParamMock,
				inputParam: updateMenuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
				mock.menu.EXPECT().Update(menuParamMock, updateMenuParamMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(mocks, tt.args)
		err := m.Update(tt.args.param, tt.args.inputParam)
		if (err != nil) != tt.wantErr {
			t.Errorf("menu.Get() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func Test_menu_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	menuParamMock := entity.MenuParam{
		ID: 1,
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		param entity.MenuParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to delete menu",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Delete(menuParamMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				param: menuParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Delete(menuParamMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt.mockFunc(mocks, tt.args)
		err := m.Delete(tt.args.param)
		if (err != nil) != tt.wantErr {
			t.Errorf("menu.Delete() error = %v, wantErr %v", err, tt.wantErr)
			return
		}
	}
}

func Test_menu_ValidateMenu(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	menuMock := mock_menu.NewMockInterface(ctrl)

	userAuthAdminMock := auth.UserAuthInfo{
		User: auth.User{
			IsAdmin: true,
			UmkmID:  1,
		},
	}

	userAuthMock := auth.UserAuthInfo{
		User: auth.User{
			IsAdmin: false,
			UmkmID:  1,
		},
	}

	menuParamMock := entity.MenuParam{
		ID: 1,
	}

	menuParam2Mock := entity.MenuParam{
		ID: 2,
	}

	menuResultMock := entity.Menu{
		UmkmID: 1,
	}

	menuResult2Mock := entity.Menu{
		UmkmID: 2,
	}

	m := menu.Init(menuMock)

	type mockFields struct {
		menu *mock_menu.MockInterface
	}

	mocks := mockFields{
		menu: menuMock,
	}

	type args struct {
		ctx    context.Context
		menuID uint
		user   auth.UserAuthInfo
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		wantErr  bool
	}{
		{
			name: "menu id 0",
			args: args{
				ctx:    context.Background(),
				menuID: 0,
				user:   userAuthMock,
			},
			mockFunc: func(mock mockFields, arg args) {},
			wantErr:  true,
		},
		{
			name: "failed to get menu",
			args: args{
				ctx:    context.Background(),
				menuID: 1,
				user:   userAuthMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(entity.Menu{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "not admin",
			args: args{
				ctx:    context.Background(),
				menuID: 2,
				user:   userAuthMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParam2Mock).Return(menuResult2Mock, nil)
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				menuID: 1,
				user:   userAuthAdminMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.menu.EXPECT().Get(menuParamMock).Return(menuResultMock, nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := m.ValidateMenu(tt.args.ctx, tt.args.menuID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("menu.ValidateMenu() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
