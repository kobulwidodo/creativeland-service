package analytic

import (
	"context"
	cartDom "go-clean/src/business/domain/cart"
	"go-clean/src/business/entity"
	"strings"
	"time"
)

type Interface interface {
	GetDashboardWidget(ctx context.Context, param entity.TransactionParam) (entity.WidgetDashboardResult, error)
}

type analytic struct {
	cart cartDom.Interface
}

func Init(cd cartDom.Interface) Interface {
	a := &analytic{
		cart: cd,
	}

	return a
}

func (a *analytic) GetDashboardWidget(ctx context.Context, param entity.TransactionParam) (entity.WidgetDashboardResult, error) {
	result := entity.WidgetDashboardResult{}

	now := time.Now()

	month := now.Format("2006-01")
	today := now.Format("2006-01-02")
	cartsMonth, err := a.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		CreatedAt: month,
	})
	if err != nil {
		return result, err
	}

	transactionMonth := make(map[uint]bool)
	revenueMonth := 0
	transactionToday := make(map[uint]bool)
	revenueToday := 0
	for _, c := range cartsMonth {
		transactionMonth[c.TransactionID] = true
		revenueMonth += c.TotalPrice
		if strings.Contains(c.CreatedAt.String(), today) {
			transactionToday[c.TransactionID] = true
			revenueToday += c.TotalPrice
		}
	}

	result.TotalMonthTransaction = len(transactionMonth)
	result.TotalMonthRevenue = revenueMonth
	result.TotalTodayTransaction = len(transactionToday)
	result.TotalTodayRevenue = revenueToday

	return result, nil
}
