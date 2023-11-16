package entity

import "gorm.io/gorm"

type Transaction struct {
	gorm.Model
	GuestID   string
	BuyerName string
	Seat      string
	Notes     string
	Price     int
}

type CreateTransactionParam struct {
	BuyerName string `binding:"required"`
	Seat      string `binding:"required"`
	Notes     string
	PaymentID int    `binding:"required"`
	Email     string `binding:"required"`
}

type TransactionParam struct {
	ID              uint     `uri:"transaction_id" json:"id"`
	Date            string   `form:"date"`
	UmkmID          uint     `uri:"umkm_id"`
	Status          string   `form:"status"`
	Statuses        []string `form:"statuses"`
	MidtransOrderID string   `form:"order_id"`
	Page            int      `form:"page" json:"-" gorm:"-"`
	Limit           int      `form:"limit" json:"-" gorm:"-"`
	Offset          int      `json:"-" gorm:"-"`
	OrderBy         string   `json:"-" gorm:"-"`
}

type TransactionDetailResponse struct {
	ID              uint        `json:"transaction_id"`
	BuyerName       string      `json:"buyer_name"`
	Seat            string      `json:"seat"`
	Notes           string      `json:"notes"`
	Price           int         `json:"price"`
	Status          string      `json:"status"`
	PaymentType     string      `json:"payment_type,omitempty"`
	MidtransOrderID string      `json:"midtrans_order_id"`
	CreatedAt       string      `json:"created_at"`
	ItemMenus       []ItemMenu  `json:"item_menus"`
	PaymentData     PaymentData `json:"payment_data"`
}

type ItemMenu struct {
	UmkmName     string `json:"umkm_name"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Price        int    `json:"price"`
	Qty          int    `json:"qty"`
	PricePerItem int    `json:"price_per_item"`
}

type SalesRecapResponse struct {
	Date        string            `json:"date"`
	GrossAmount int               `json:"gross_amount"`
	NetAmount   int               `json:"net_amount"`
	UmkmDetail  []UmkmDetailRecap `json:"umkm_detail"`
}

type UmkmDetailRecap struct {
	ID          uint   `json:"id"`
	UmkmName    string `json:"umkm_name"`
	GrossAmount int    `json:"gross_amount"`
	NetAmount   int    `json:"net_amount"`
	TotalOrder  int    `json:"total_order"`
}

type KeyUmkmDetailRecap struct {
	ID          uint
	CreatedDate string
}
