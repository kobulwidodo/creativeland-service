package umkm_test

import (
	"context"
	mock_umkm "go-clean/src/business/domain/mock/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/umkm"
	"go-clean/src/lib/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_umkm_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	umkmMock := mock_umkm.NewMockInterface(ctrl)

	paramMock := entity.CreateUmkmParam{
		Name:   "umkm",
		Slogan: "slogan",
	}

	umkmResultMock := entity.Umkm{
		Model: gorm.Model{
			ID: 1,
		},
		Name:   "umkm",
		Slogan: "slogan",
	}

	u := umkm.Init(umkmMock)

	type mockfields struct {
		umkm *mock_umkm.MockInterface
	}

	mocks := mockfields{
		umkm: umkmMock,
	}

	type args struct {
		params entity.CreateUmkmParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.Umkm
		wantErr  bool
	}{
		{
			name: "failed to create umkm",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Create(gomock.Any()).Return(entity.Umkm{}, assert.AnError)
			},
			args: args{
				params: paramMock,
			},
			want:    entity.Umkm{},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Create(gomock.Any()).Return(umkmResultMock, nil)
			},
			args: args{
				params: paramMock,
			},
			want:    umkmResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Create(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_umkm_GetList(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	umkmMock := mock_umkm.NewMockInterface(ctrl)

	paramMock := entity.UmkmParam{
		ID: 1,
	}

	umkmResultMock := []entity.Umkm{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name:   "umkm",
			Slogan: "slogan",
		},
	}

	u := umkm.Init(umkmMock)

	type mockfields struct {
		umkm *mock_umkm.MockInterface
	}

	mocks := mockfields{
		umkm: umkmMock,
	}

	type args struct {
		params entity.UmkmParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     []entity.Umkm
		wantErr  bool
	}{
		{
			name: "failed to get umkm list",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().GetList(gomock.Any()).Return([]entity.Umkm{}, assert.AnError)
			},
			args: args{
				params: paramMock,
			},
			want:    []entity.Umkm{},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().GetList(paramMock).Return(umkmResultMock, nil)
			},
			args: args{
				params: paramMock,
			},
			want:    umkmResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.GetList(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.GetList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_umkm_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	umkmMock := mock_umkm.NewMockInterface(ctrl)

	paramMock := entity.UmkmParam{
		ID: 1,
	}

	umkmResultMock := entity.Umkm{
		Model: gorm.Model{
			ID: 1,
		},
		Name:   "umkm",
		Slogan: "slogan",
	}

	u := umkm.Init(umkmMock)

	type mockfields struct {
		umkm *mock_umkm.MockInterface
	}

	mocks := mockfields{
		umkm: umkmMock,
	}

	type args struct {
		params entity.UmkmParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.Umkm
		wantErr  bool
	}{
		{
			name: "failed to get umkm list",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Get(gomock.Any()).Return(entity.Umkm{}, assert.AnError)
			},
			args: args{
				params: paramMock,
			},
			want:    entity.Umkm{},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Get(paramMock).Return(umkmResultMock, nil)
			},
			args: args{
				params: paramMock,
			},
			want:    umkmResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_umkm_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	umkmMock := mock_umkm.NewMockInterface(ctrl)

	paramMock := entity.UmkmParam{
		ID: 1,
	}

	updateParamMock := entity.UpdateUmkmParam{
		Name: "new name",
	}

	umkmResultMock := entity.Umkm{
		Model: gorm.Model{
			ID: 1,
		},
	}

	u := umkm.Init(umkmMock)

	type mockfields struct {
		umkm *mock_umkm.MockInterface
	}

	mocks := mockfields{
		umkm: umkmMock,
	}

	type args struct {
		params      entity.UmkmParam
		updateParam entity.UpdateUmkmParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		wantErr  bool
	}{
		{
			name: "failed to get umkm",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Get(paramMock).Return(entity.Umkm{}, assert.AnError)
			},
			args: args{
				params: paramMock,
			},
			wantErr: true,
		},
		{
			name: "failed to update umkm",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Get(paramMock).Return(umkmResultMock, nil)
				mock.umkm.EXPECT().Update(paramMock, updateParamMock).Return(assert.AnError)
			},
			args: args{
				params:      paramMock,
				updateParam: updateParamMock,
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Get(paramMock).Return(umkmResultMock, nil)
				mock.umkm.EXPECT().Update(paramMock, updateParamMock).Return(nil)
			},
			args: args{
				params:      paramMock,
				updateParam: updateParamMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := u.Update(tt.args.params, tt.args.updateParam)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_umkm_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	umkmMock := mock_umkm.NewMockInterface(ctrl)

	paramMock := entity.UmkmParam{
		ID: 1,
	}

	u := umkm.Init(umkmMock)

	type mockfields struct {
		umkm *mock_umkm.MockInterface
	}

	mocks := mockfields{
		umkm: umkmMock,
	}

	type args struct {
		params entity.UmkmParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		wantErr  bool
	}{
		{
			name: "failed to delete umkm",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Delete(paramMock).Return(assert.AnError)
			},
			args: args{
				params: paramMock,
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.umkm.EXPECT().Delete(paramMock).Return(nil)
			},
			args: args{
				params: paramMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := u.Delete(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_umkm_ValidateUmkm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

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

	u := umkm.Init(nil)

	type args struct {
		ctx    context.Context
		umkmId uint
		user   auth.UserAuthInfo
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "umkm id 0",
			args: args{
				ctx:    context.Background(),
				umkmId: 0,
				user:   userAuthMock,
			},
			wantErr: true,
		},
		{
			name: "not admin",
			args: args{
				ctx:    context.Background(),
				umkmId: 2,
				user:   userAuthMock,
			},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:    context.Background(),
				umkmId: 1,
				user:   userAuthAdminMock,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := u.ValidateUmkm(tt.args.ctx, tt.args.umkmId, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("umkm.ValidateUmkm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
