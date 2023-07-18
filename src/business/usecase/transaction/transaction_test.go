package transaction_test

import (
	"context"
	"encoding/json"
	mock_cart "go-clean/src/business/domain/mock/cart"
	mock_menu "go-clean/src/business/domain/mock/menu"
	mock_midtrans "go-clean/src/business/domain/mock/midtrans"
	mock_midtrans_transaction "go-clean/src/business/domain/mock/midtrans_transaction"
	mock_transaction "go-clean/src/business/domain/mock/transaction"
	mock_umkm "go-clean/src/business/domain/mock/umkm"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/transaction"
	"go-clean/src/lib/auth"
	"go-clean/src/lib/midtrans"
	mock_auth "go-clean/src/lib/tests/mock/auth"
	"testing"

	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_transaction_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	authMock := mock_auth.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)
	midtransMock := mock_midtrans.NewMockInterface(ctrl)
	transactionMock := mock_transaction.NewMockInterface(ctrl)
	midtransTransactionMock := mock_midtrans_transaction.NewMockInterface(ctrl)

	tr := transaction.Init(authMock, transactionMock, cartMock, menuMock, nil, midtransMock, midtransTransactionMock)

	userAuthMock := auth.UserAuthInfo{
		User: auth.User{
			GuestID: "1",
		},
	}

	paramsMock := entity.CreateTransactionParam{
		BuyerName: "mail",
		Seat:      "a1",
		Notes:     "-",
		PaymentID: 1,
		Email:     "mail@gmail.com",
	}

	paramsMockUndifinedMock := entity.CreateTransactionParam{
		BuyerName: "mail",
		Seat:      "a1",
		Notes:     "-",
		PaymentID: 999,
		Email:     "mail@gmail.com",
	}

	cartParamMock := entity.CartParam{
		Status:  entity.StatusInCart,
		GuestID: "1",
	}

	cartResultMock := []entity.Cart{
		{
			Model: gorm.Model{
				ID: 1,
			},
			MenuID:       1,
			TotalPrice:   10000,
			PricePerItem: 10000,
			Amount:       1,
		},
	}

	menuResultMock := []entity.Menu{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "menu 1",
		},
	}

	newTransactionMock := entity.Transaction{
		GuestID:   "1",
		BuyerName: "mail",
		Seat:      "a1",
		Notes:     "-",
		Price:     10000,
	}

	transactionResultMock := entity.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		GuestID:   "1",
		BuyerName: "mail",
		Seat:      "a1",
		Notes:     "-",
		Price:     10000,
	}

	midtransCreateParamMock := midtrans.CreateOrderParam{
		OrderID:     1,
		PaymentID:   1,
		GrossAmount: 10000,
		ItemsDetails: []midtrans.ItemsDetails{
			{
				ID:    "1",
				Price: 10000,
				Qty:   1,
				Name:  "menu 1",
			},
		},
		CustomerDetails: midtrans.CustomerDetails{
			Name:  "mail",
			Email: "mail@gmail.com",
		},
	}

	midtransCreateParamUndifinedMock := midtrans.CreateOrderParam{
		OrderID:     1,
		PaymentID:   999,
		GrossAmount: 10000,
		ItemsDetails: []midtrans.ItemsDetails{
			{
				ID:    "1",
				Price: 10000,
				Qty:   1,
				Name:  "menu 1",
			},
		},
		CustomerDetails: midtrans.CustomerDetails{
			Name:  "mail",
			Email: "mail@gmail.com",
		},
	}

	midtransResultMock := &coreapi.ChargeResponse{
		TransactionID: "1",
		OrderID:       "1",
		Actions: []coreapi.Action{
			{
				URL: "url 1",
			},
			{
				URL: "url 2",
			},
		},
	}

	paymentData, _ := json.Marshal(entity.PaymentData{
		Key: "url 2",
		Qr:  "url 1",
	})

	newMidtransTransactionMock := entity.MidtransTransaction{
		TransactionID: 1,
		MidtransID:    "1",
		OrderID:       "1",
		PaymentType:   1,
		GrossAmount:   10000,
		Status:        "pending",
		PaymentData:   string(paymentData),
	}

	selectParamCartMock := entity.CartParam{
		Status:  entity.StatusInCart,
		GuestID: "1",
	}

	updateParamCartMock := entity.UpdateCartParam{
		Status:        entity.StatusUnpaid,
		TransactionID: 1,
	}

	type mockfields struct {
		auth                 *mock_auth.MockInterface
		cart                 *mock_cart.MockInterface
		menu                 *mock_menu.MockInterface
		midtrans             *mock_midtrans.MockInterface
		transaction          *mock_transaction.MockInterface
		midtrans_transaction *mock_midtrans_transaction.MockInterface
	}

	mocks := mockfields{
		auth:                 authMock,
		cart:                 cartMock,
		menu:                 menuMock,
		midtrans:             midtransMock,
		transaction:          transactionMock,
		midtrans_transaction: midtransTransactionMock,
	}

	type args struct {
		ctx   context.Context
		param entity.CreateTransactionParam
	}

	tests := []struct {
		name     string
		mockFunc func(mock mockfields, arg args)
		args     args
		want     uint
		wantErr  bool
	}{
		{
			name: "failed to get auth user",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed get cart list",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "cart empty",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to get menus list",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return([]entity.Menu{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to create transaction",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to create midtrans",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(nil, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to update cart",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to get payment data",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamUndifinedMock).Return(midtransResultMock, nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMockUndifinedMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "failed to create midtrans transaction",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(nil)
				mock.midtrans_transaction.EXPECT().Create(newMidtransTransactionMock).Return(entity.MidtransTransaction{}, assert.AnError)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "all success",
			mockFunc: func(mock mockfields, arg args) {
				mock.auth.EXPECT().GetUserAuthInfo(context.Background()).Return(userAuthMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
				mock.transaction.EXPECT().Create(newTransactionMock).Return(transactionResultMock, nil)
				mock.midtrans.EXPECT().Create(midtransCreateParamMock).Return(midtransResultMock, nil)
				mock.cart.EXPECT().Update(selectParamCartMock, updateParamCartMock).Return(nil)
				mock.midtrans_transaction.EXPECT().Create(newMidtransTransactionMock).Return(entity.MidtransTransaction{}, nil)
			},
			args: args{
				ctx:   context.Background(),
				param: paramsMock,
			},
			want:    1,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := tr.Create(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transaction_GetTransactionListByUmkm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartMock := mock_cart.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)
	transactionMock := mock_transaction.NewMockInterface(ctrl)
	midtransTransactionMock := mock_midtrans_transaction.NewMockInterface(ctrl)

	tr := transaction.Init(nil, transactionMock, cartMock, menuMock, nil, nil, midtransTransactionMock)

	transactionParamMock := entity.TransactionParam{
		UmkmID:          1,
		Status:          entity.StatusPaid,
		MidtransOrderID: "",
	}

	cartParamMock := entity.CartParam{
		UmkmID: 1,
		Status: entity.StatusPaid,
	}

	cartResultMock := []entity.Cart{
		{
			TransactionID: 1,
			Status:        entity.StatusDone,
			TotalPrice:    10000,
			Amount:        1,
			PricePerItem:  10000,
		},
	}

	menuParamMock := entity.MenuParam{
		UmkmID: 1,
	}

	menuResultMock := []entity.Menu{
		{
			Name: "menu 1",
		},
	}

	transactionResultMock := []entity.Transaction{
		{
			Model: gorm.Model{
				ID: 1,
			},
			BuyerName: "mail",
			Seat:      "a1",
			Notes:     "-",
			Price:     10000,
		},
	}

	midtransTransactionParamMock := entity.MidtransTransactionParam{
		OrderID: "",
	}

	midtransTransactionResultMock := []entity.MidtransTransaction{
		{
			OrderID:       "1",
			TransactionID: 1,
		},
	}

	resultMock := []entity.TransactionDetailResponse{
		{
			ID:              1,
			BuyerName:       "mail",
			Seat:            "a1",
			Notes:           "-",
			Price:           10000,
			Status:          entity.StatusDone,
			MidtransOrderID: "1",
			ItemMenus: []entity.ItemMenu{
				{
					Name:         "menu 1",
					Price:        10000,
					Qty:          1,
					PricePerItem: 10000,
				},
			},
		},
	}

	type mockfields struct {
		cart                 *mock_cart.MockInterface
		menu                 *mock_menu.MockInterface
		transaction          *mock_transaction.MockInterface
		midtrans_transaction *mock_midtrans_transaction.MockInterface
	}

	mocks := mockfields{
		cart:                 cartMock,
		menu:                 menuMock,
		transaction:          transactionMock,
		midtrans_transaction: midtransTransactionMock,
	}

	type args struct {
		ctx   context.Context
		param entity.TransactionParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockfields, arg args)
		want     []entity.TransactionDetailResponse
		wantErr  bool
	}{
		{
			name: "failed to get cart list",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    []entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "cart empty",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, nil)
			},
			want:    []entity.TransactionDetailResponse{},
			wantErr: false,
		},
		{
			name: "failed to get menu list",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetAll(menuParamMock).Return([]entity.Menu{}, assert.AnError)
			},
			want:    []entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed to get transaction list by id",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetAll(menuParamMock).Return(menuResultMock, nil)
				mock.transaction.EXPECT().GetListByIDs([]uint{1}).Return([]entity.Transaction{}, assert.AnError)
			},
			want:    []entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed to get midtrans transaction list by id",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetAll(menuParamMock).Return(menuResultMock, nil)
				mock.transaction.EXPECT().GetListByIDs([]uint{1}).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().GetListByTrxIDs([]uint{1}, midtransTransactionParamMock).Return([]entity.MidtransTransaction{}, assert.AnError)
			},
			want:    []entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "success",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.menu.EXPECT().GetAll(menuParamMock).Return(menuResultMock, nil)
				mock.transaction.EXPECT().GetListByIDs([]uint{1}).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().GetListByTrxIDs([]uint{1}, midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := tr.GetTransactionListByUmkm(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.GetTransactionListByUmkm() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transaction_GetOrderDetail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	transactionMock := mock_transaction.NewMockInterface(ctrl)
	midtransTransactionMock := mock_midtrans_transaction.NewMockInterface(ctrl)
	cartMock := mock_cart.NewMockInterface(ctrl)
	umkmMock := mock_umkm.NewMockInterface(ctrl)
	menuMock := mock_menu.NewMockInterface(ctrl)

	transactionParamMock := entity.TransactionParam{
		ID: 1,
	}

	transactionResultMock := entity.Transaction{
		Model: gorm.Model{
			ID: 1,
		},
		BuyerName: "mail",
		Seat:      "a1",
		Notes:     "-",
		Price:     10000,
	}

	midtransTransactionParamMock := entity.MidtransTransactionParam{
		TransactionID: 1,
	}

	midtransTransactionResultMock := entity.MidtransTransaction{
		Status:      entity.StatusPaid,
		PaymentType: midtrans.GopayPayment,
	}

	cartParamMock := entity.CartParam{
		TransactionID: 1,
	}

	cartResultMock := []entity.Cart{
		{
			MenuID:       1,
			UmkmID:       1,
			TotalPrice:   10000,
			Amount:       1,
			PricePerItem: 10000,
		},
	}

	umkmResultMock := []entity.Umkm{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "umkm 1",
		},
	}

	menuResultMock := []entity.Menu{
		{
			Model: gorm.Model{
				ID: 1,
			},
			Name: "menu 1",
		},
	}

	resultMock := entity.TransactionDetailResponse{
		ID:          1,
		BuyerName:   "mail",
		Seat:        "a1",
		Notes:       "-",
		Price:       10000,
		Status:      entity.StatusPaid,
		PaymentType: "Gopay",
		ItemMenus: []entity.ItemMenu{
			{
				UmkmName:     "umkm 1",
				Name:         "menu 1",
				Price:        10000,
				Qty:          1,
				PricePerItem: 10000,
			},
		},
	}

	tr := transaction.Init(nil, transactionMock, cartMock, menuMock, umkmMock, nil, midtransTransactionMock)

	type mockfields struct {
		cart                 *mock_cart.MockInterface
		menu                 *mock_menu.MockInterface
		umkm                 *mock_umkm.MockInterface
		transaction          *mock_transaction.MockInterface
		midtrans_transaction *mock_midtrans_transaction.MockInterface
	}

	mocks := mockfields{
		cart:                 cartMock,
		menu:                 menuMock,
		umkm:                 umkmMock,
		transaction:          transactionMock,
		midtrans_transaction: midtransTransactionMock,
	}

	type args struct {
		ctx   context.Context
		param entity.TransactionParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockfields, arg args)
		want     entity.TransactionDetailResponse
		wantErr  bool
	}{
		{
			name: "failed get transaction",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(entity.Transaction{}, assert.AnError)
			},
			want:    entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed get midtrans transaction",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(entity.MidtransTransaction{}, assert.AnError)
			},
			want:    entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed get cart list",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed get umkm list",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.umkm.EXPECT().GetList(entity.UmkmParam{}).Return([]entity.Umkm{}, assert.AnError)
			},
			want:    entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "failed get menu list by id",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.umkm.EXPECT().GetList(entity.UmkmParam{}).Return(umkmResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return([]entity.Menu{}, assert.AnError)
			},
			want:    entity.TransactionDetailResponse{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.transaction.EXPECT().Get(transactionParamMock).Return(transactionResultMock, nil)
				mock.midtrans_transaction.EXPECT().Get(midtransTransactionParamMock).Return(midtransTransactionResultMock, nil)
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.umkm.EXPECT().GetList(entity.UmkmParam{}).Return(umkmResultMock, nil)
				mock.menu.EXPECT().GetListInByID([]int64{1}).Return(menuResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := tr.GetOrderDetail(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.GetOrderDetail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_transaction_CompleteOrder(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cartMock := mock_cart.NewMockInterface(ctrl)

	transactionParamMock := entity.TransactionParam{
		ID:     1,
		UmkmID: 1,
	}

	cartParamMock := entity.CartParam{
		TransactionID: 1,
		UmkmID:        1,
		Status:        entity.StatusPaid,
	}

	cartResultMock := []entity.Cart{
		{
			Model: gorm.Model{
				ID: 1,
			},
		},
	}

	updateCartParamMock := entity.UpdateCartParam{
		Status: entity.StatusDone,
	}

	tr := transaction.Init(nil, nil, cartMock, nil, nil, nil, nil)

	type mockfields struct {
		cart *mock_cart.MockInterface
	}

	mocks := mockfields{
		cart: cartMock,
	}

	type args struct {
		ctx   context.Context
		param entity.TransactionParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockfields, arg args)
		wantErr  bool
	}{
		{
			name: "failed to get cart list",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "failed to update by ids",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().UpdatesByIDs([]uint{1}, updateCartParamMock).Return(assert.AnError)
			},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:   context.Background(),
				param: transactionParamMock,
			},
			mockFunc: func(mock mockfields, arg args) {
				mock.cart.EXPECT().GetList(cartParamMock).Return(cartResultMock, nil)
				mock.cart.EXPECT().UpdatesByIDs([]uint{1}, updateCartParamMock).Return(nil)
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			err := tr.CompleteOrder(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("transaction.CompleteOrder() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
