package entity

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string
	Description string
	Price       int
	UmkmID      uint
}

type MenuParam struct {
	ID     uint   `uri:"menu_id" json:"id"`
	Name   string `form:"name" json:"name"`
	UmkmID uint   `form:"umkm_id" json:"umkm_id"`
}

type CreateMenuParam struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Price       int    `binding:"required"`
	UmkmID      uint   `binding:"required"`
}

type UpdateMenuParam struct {
	Name        string
	Description string
	Price       int
}
