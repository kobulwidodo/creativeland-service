package entity

type AnalyticParam struct {
	UmkmID uint `uri:"umkm_id"`
}

type WidgetDashboardResult struct {
	TotalTodayTransaction     int
	TotalTodayRevenue         int
	TotalYesterdayTransaction int
	TotalYesterdayRevenue     int
	TotalMonthTransaction     int
	TotalMonthRevenue         int
	TotalLastMonthTransaction int
	TotalLastMonthRevenue     int
}
