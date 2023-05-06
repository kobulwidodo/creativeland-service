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
