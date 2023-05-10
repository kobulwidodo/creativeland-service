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
	Notes     string `binding:"required"`
	PaymentID int    `binding:"required"`
	Email     string `binding:"required"`
}

type TransactionParam struct {
	ID uint `uri:"transaction_id" json:"id"`
}

type TransactionDetailResponse struct {
	ID          uint       `json:"order_id"`
	BuyerName   string     `json:"buyer_name"`
	Seat        string     `json:"seat"`
	Notes       string     `json:"notes"`
	Price       int        `json:"price"`
	Status      string     `json:"status"`
	PaymentType string     `json:"payment_type"`
	ItemMenus   []ItemMenu `json:"item_menus"`
}

type ItemMenu struct {
	Name         string `json:"name"`
	Price        int    `json:"price"`
	Qty          int    `json:"qty"`
	PricePerItem int    `json:"price_per_item"`
}
