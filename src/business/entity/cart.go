package entity

import (
	"gorm.io/gorm"
)

const (
	StatusActive   = "active"
	StatusInactive = "inactive"
)

type Cart struct {
	gorm.Model
	UmkmID        uint
	MenuID        uint
	TransactionID uint
	Status        string
	GuestID       string
	Amount        int
	TotalPrice    int
	PricePerItem  int
}

type CartParam struct {
	ID      uint `uri:"cart_id" json:"id"`
	Status  string
	GuestID string
}

type CreateCartParam struct {
	UmkmID uint `binding:"required"`
	MenuID uint `binding:"required"`
	Amount int  `binding:"required"`
}

type UpdateCartParam struct {
	TransactionID uint
	Status        string
}
