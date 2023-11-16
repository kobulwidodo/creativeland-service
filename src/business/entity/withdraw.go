package entity

import "gorm.io/gorm"

type Withdraw struct {
	gorm.Model
	Date     string
	Amount   int
	UmkmID   uint
	Status   string
	Method   string
	UmkmName string `grom:"-:all"`
}

type WithdrawParam struct {
	ID      uint   `json:"withdraw_id" uri:"withdraw_id"`
	Date    string `form:"date"`
	UmkmID  uint   `form:"umkm_id"`
	Limit   int    `form:"limit" json:"-" gorm:"-"`
	Offset  int    `json:"-" gorm:"-"`
	OrderBy string `json:"-" gorm:"-"`
	Page    int    `form:"page" json:"-" gorm:"-"`
}

type CreateWithdrawParam struct {
	Date   string `binding:"required"`
	Amount int    `binding:"required"`
	UmkmID uint   `binding:"required"`
	Status string `binding:"required"`
	Method string `binding:"required"`
}

type UpdateWithdrawParam struct {
	Amount int
	Status string
}
