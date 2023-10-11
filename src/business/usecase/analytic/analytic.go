package analytic

import (
	"context"
	cartDom "go-clean/src/business/domain/cart"
	"go-clean/src/business/entity"
	"strings"
	"time"
)

type Interface interface {
	GetDashboardWidget(ctx context.Context, param entity.AnalyticParam) (entity.WidgetDashboardResult, error)
	GetAllDashboardWidget(ctx context.Context) (entity.WidgetDashboardResult, error)
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

func (a *analytic) GetDashboardWidget(ctx context.Context, param entity.AnalyticParam) (entity.WidgetDashboardResult, error) {
	result := entity.WidgetDashboardResult{}

	now := time.Now()

	month := now.Format("2006-01")
	calcLastMonth := now.AddDate(0, -1, 0)
	lastMonth := calcLastMonth.Format("2006-01")
	today := now.Format("2006-01-02")
	calcYesterday := now.AddDate(0, 0, -1)
	yesterday := calcYesterday.Format("2006-01-02")

	cartsMonth, err := a.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		UmkmID:    param.UmkmID,
		CreatedAt: month,
	})
	if err != nil {
		return result, err
	}

	cartsLastsMonth, err := a.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		UmkmID:    param.UmkmID,
		CreatedAt: lastMonth,
	})
	if err != nil {
		return result, err
	}

	mergedCart := append(cartsMonth, cartsLastsMonth...)

	transactionMonth := make(map[uint]bool)
	revenueMonth := 0
	transactionLastMonth := make(map[uint]bool)
	revenueLastMonth := 0
	transactionToday := make(map[uint]bool)
	revenueToday := 0
	transactionYesterday := make(map[uint]bool)
	revenueYesterday := 0
	for _, c := range mergedCart {
		if strings.Contains(c.CreatedAt.String(), month) {
			transactionMonth[c.TransactionID] = true
			revenueMonth += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), today) {
			transactionToday[c.TransactionID] = true
			revenueToday += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), lastMonth) {
			transactionLastMonth[c.TransactionID] = true
			revenueLastMonth += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), yesterday) {
			transactionYesterday[c.TransactionID] = true
			revenueYesterday += c.TotalPrice
		}
	}

	result.TotalMonthTransaction = len(transactionMonth)
	result.TotalMonthRevenue = revenueMonth
	result.TotalLastMonthTransaction = len(transactionLastMonth)
	result.TotalLastMonthRevenue = revenueLastMonth
	result.TotalTodayTransaction = len(transactionToday)
	result.TotalTodayRevenue = revenueToday
	result.TotalYesterdayTransaction = len(transactionYesterday)
	result.TotalYesterdayRevenue = revenueYesterday

	return result, nil
}

func (a *analytic) GetAllDashboardWidget(ctx context.Context) (entity.WidgetDashboardResult, error) {
	result := entity.WidgetDashboardResult{}

	now := time.Now()

	month := now.Format("2006-01")
	calcLastMonth := now.AddDate(0, -1, 0)
	lastMonth := calcLastMonth.Format("2006-01")
	today := now.Format("2006-01-02")
	calcYesterday := now.AddDate(0, 0, -1)
	yesterday := calcYesterday.Format("2006-01-02")

	cartsMonth, err := a.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		CreatedAt: month,
	})
	if err != nil {
		return result, err
	}

	cartsLastsMonth, err := a.cart.GetList(entity.CartParam{
		Status:    entity.StatusDone,
		CreatedAt: lastMonth,
	})
	if err != nil {
		return result, err
	}

	mergedCart := append(cartsMonth, cartsLastsMonth...)

	transactionMonth := make(map[uint]bool)
	revenueMonth := 0
	transactionLastMonth := make(map[uint]bool)
	revenueLastMonth := 0
	transactionToday := make(map[uint]bool)
	revenueToday := 0
	transactionYesterday := make(map[uint]bool)
	revenueYesterday := 0
	for _, c := range mergedCart {
		if strings.Contains(c.CreatedAt.String(), month) {
			transactionMonth[c.TransactionID] = true
			revenueMonth += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), today) {
			transactionToday[c.TransactionID] = true
			revenueToday += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), lastMonth) {
			transactionLastMonth[c.TransactionID] = true
			revenueLastMonth += c.TotalPrice
		}
		if strings.Contains(c.CreatedAt.String(), yesterday) {
			transactionYesterday[c.TransactionID] = true
			revenueYesterday += c.TotalPrice
		}
	}

	result.TotalMonthTransaction = len(transactionMonth)
	result.TotalMonthRevenue = revenueMonth
	result.TotalLastMonthTransaction = len(transactionLastMonth)
	result.TotalLastMonthRevenue = revenueLastMonth
	result.TotalTodayTransaction = len(transactionToday)
	result.TotalTodayRevenue = revenueToday
	result.TotalYesterdayTransaction = len(transactionYesterday)
	result.TotalYesterdayRevenue = revenueYesterday

	return result, nil
}
