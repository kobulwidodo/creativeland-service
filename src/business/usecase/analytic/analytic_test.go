package analytic_test

import (
	"context"
	mock_cart "go-clean/src/business/domain/mock/cart"
	"go-clean/src/business/entity"
	"go-clean/src/business/usecase/analytic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func Test_analytic_GetDashboardWidget(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	now := time.Now()
	yesterday := now.AddDate(0, 0, -1)
	lastMonth := now.AddDate(0, -1, 0)
	thisMonthFormatted := now.Format("2006-01")
	lastMonthFormatted := lastMonth.Format("2006-01")

	cartMock := mock_cart.NewMockInterface(ctrl)

	analyticParamMock := entity.AnalyticParam{
		UmkmID: 1,
	}

	cartThisMonthParamMock := entity.CartParam{
		Status:    entity.StatusDone,
		UmkmID:    1,
		CreatedAt: thisMonthFormatted,
	}

	cartLastMonthParamMock := entity.CartParam{
		Status:    entity.StatusDone,
		UmkmID:    1,
		CreatedAt: lastMonthFormatted,
	}

	cartsThisMonthResultMock := []entity.Cart{
		{
			TransactionID: 1,
			TotalPrice:    10000,
			Model: gorm.Model{
				CreatedAt: now,
			},
		},
		{
			TransactionID: 2,
			TotalPrice:    20000,
			Model: gorm.Model{
				CreatedAt: yesterday,
			},
		},
	}

	cartsLastMonthResultMock := []entity.Cart{
		{
			TransactionID: 3,
			TotalPrice:    30000,
			Model: gorm.Model{
				CreatedAt: lastMonth,
			},
		},
	}

	resultMock := entity.WidgetDashboardResult{
		TotalMonthTransaction:     2,
		TotalMonthRevenue:         30000,
		TotalLastMonthTransaction: 1,
		TotalLastMonthRevenue:     30000,
		TotalTodayTransaction:     1,
		TotalTodayRevenue:         10000,
		TotalYesterdayTransaction: 1,
		TotalYesterdayRevenue:     20000,
	}

	a := analytic.Init(cartMock)

	type mockFields struct {
		cart *mock_cart.MockInterface
	}

	mocks := mockFields{
		cart: cartMock,
	}

	type args struct {
		ctx   context.Context
		param entity.AnalyticParam
	}

	tests := []struct {
		name     string
		args     args
		mockFunc func(mock mockFields, arg args)
		want     entity.WidgetDashboardResult
		wantErr  bool
	}{
		{
			name: "failed to get cart list",
			args: args{
				ctx:   context.Background(),
				param: analyticParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().GetList(cartThisMonthParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    entity.WidgetDashboardResult{},
			wantErr: true,
		},
		{
			name: "failed to get last month cart list",
			args: args{
				ctx:   context.Background(),
				param: analyticParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().GetList(cartThisMonthParamMock).Return(cartsThisMonthResultMock, nil)
				mock.cart.EXPECT().GetList(cartLastMonthParamMock).Return([]entity.Cart{}, assert.AnError)
			},
			want:    entity.WidgetDashboardResult{},
			wantErr: true,
		},
		{
			name: "all success",
			args: args{
				ctx:   context.Background(),
				param: analyticParamMock,
			},
			mockFunc: func(mock mockFields, arg args) {
				mock.cart.EXPECT().GetList(cartThisMonthParamMock).Return(cartsThisMonthResultMock, nil)
				mock.cart.EXPECT().GetList(cartLastMonthParamMock).Return(cartsLastMonthResultMock, nil)
			},
			want:    resultMock,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc(mocks, tt.args)
			got, err := a.GetDashboardWidget(tt.args.ctx, tt.args.param)
			if (err != nil) != tt.wantErr {
				t.Errorf("analytic.GetDashboardWidget() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
