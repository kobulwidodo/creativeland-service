package entity

import (
	"gorm.io/gorm"
)

const (
	StatusInCart = "in_cart"
	StatusUnpaid = "unpaid"
	StatusPaid   = "paid"
	StatusDone   = "done"
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
	Menu          Menu `grom:"-:all"`
	Umkm          Umkm `grom:"-:all"`
}

type CartParam struct {
	ID            uint `uri:"cart_id" json:"id"`
	TransactionID uint
	Status        string
	UmkmID        uint
	MenuID        uint
	GuestID       string
}

type CreateCartParam struct {
	UmkmID uint `binding:"required"`
	MenuID uint `binding:"required"`
	Amount int  `binding:"required"`
}

type UpdateCartParam struct {
	TransactionID uint
	Status        string
	TotalPrice    int
	Amount        int
}
