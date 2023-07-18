package user_test

import (
	"context"
	mock_user "go-clean/src/business/domain/mock/user"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/user"
	"go-clean/src/lib/auth"
	mock_auth "go-clean/src/lib/tests/mock/auth"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Test_user_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)

	mockParams := entity.CreateUserParam{
		Username: "username",
		Password: "password",
		Nama:     "nama",
		UmkmID:   1,
	}

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "username",
		Password: string(hashPass),
	}

	u := user.Init(userMock, nil, nil)

	type mockfields struct {
		user *mock_user.MockInterface
	}

	mocks := mockfields{
		user: userMock,
	}

	type args struct {
		params entity.CreateUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to create user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Username: "username",
			},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Create(gomock.Any()).Return(mockUserResult, nil)
			},
			args: args{
				params: mockParams,
			},
			want: entity.User{
				Username: "username",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Create(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want.Username, got.Username)
		})
	}
}

func Test_user_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)
	authMock := mock_auth.NewMockInterface(ctrl)
	hashPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.MinCost)
	tokenMock := "token"

	mockParams := entity.LoginUserParam{
		Username: "username",
		Password: "password",
	}

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "username",
		Password: string(hashPass),
	}

	u := user.Init(userMock, authMock, nil)

	type mockfields struct {
		user *mock_user.MockInterface
		auth *mock_auth.MockInterface
	}

	mocks := mockfields{
		user: userMock,
		auth: authMock,
	}

	type args struct {
		params entity.LoginUserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     string
		wantErr  bool
	}{
		{
			name: "failed to find user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(entity.User{}, assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "user not found",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(entity.User{}, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "password incorrect",
			mockFunc: func(mock mockfields, arg args) {
				mockUserResultWithWrongPassword := mockUserResult
				mockUserResultWithWrongPassword.Password = "wrongPassword"
				mock.user.EXPECT().Get(gomock.Any()).Return(mockUserResultWithWrongPassword, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed to generate token",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return("", assert.AnError)
			},
			args: args{
				params: mockParams,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(mockUserResult, nil)
				mock.auth.EXPECT().GenerateToken(gomock.Any()).Return(tokenMock, nil)
			},
			args: args{
				params: mockParams,
			},
			want:    tokenMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Login(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userMock := mock_user.NewMockInterface(ctrl)

	paramsMock := entity.UserParam{
		ID: 1,
	}

	userResultMock := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
		Username: "username",
	}

	u := user.Init(userMock, nil, nil)

	type mockfields struct {
		user *mock_user.MockInterface
	}

	mocks := mockfields{
		user: userMock,
	}

	type args struct {
		params entity.UserParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to get user",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(entity.User{}, assert.AnError)
			},
			args: args{
				params: paramsMock,
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields, arg args) {
				mock.user.EXPECT().Get(gomock.Any()).Return(userResultMock, nil)
			},
			args: args{
				params: paramsMock,
			},
			want:    userResultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Get(tt.args.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_GenerateGuestToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	tokenMock := "token"

	u := user.Init(nil, authMock, nil)

	type mockfields struct {
		auth *mock_auth.MockInterface
	}

	mocks := mockfields{
		auth: authMock,
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields)
		want     string
		wantErr  bool
	}{
		{
			name: "failed to generate guest token",
			mockFunc: func(mock mockfields) {
				mock.auth.EXPECT().GenerateGuestToken().Return("", assert.AnError)
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "success",
			mockFunc: func(mock mockfields) {
				mock.auth.EXPECT().GenerateGuestToken().Return(tokenMock, nil)
			},
			want:    tokenMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks)
			got, err := u.GenerateGuestToken()
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_user_Me(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	userMock := mock_user.NewMockInterface(ctrl)

	u := user.Init(userMock, authMock, nil)

	mockAuthUserInfo := auth.UserAuthInfo{
		User: auth.User{
			ID: 1,
		},
	}

	mockUserParam := entity.UserParam{
		ID: 1,
	}

	mockUserResult := entity.User{
		Model: gorm.Model{
			ID: 1,
		},
	}

	type mockfields struct {
		user *mock_user.MockInterface
		auth *mock_auth.MockInterface
	}

	mocks := mockfields{
		user: userMock,
		auth: authMock,
	}

	type args struct {
		ctx context.Context
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockfields, arg args)
		want     entity.User
		wantErr  bool
	}{
		{
			name: "failed to get user auth",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(auth.UserAuthInfo{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "failed to get user",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockAuthUserInfo, nil)
				mock.user.EXPECT().Get(mockUserParam).Return(entity.User{}, assert.AnError)
			},
			want:    entity.User{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(gomock.Any()).Return(mockAuthUserInfo, nil)
				mock.user.EXPECT().Get(mockUserParam).Return(mockUserResult, nil)
			},
			want:    mockUserResult,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := u.Me(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("user.Me() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
