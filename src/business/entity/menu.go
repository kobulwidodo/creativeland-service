package entity

import "gorm.io/gorm"

type Menu struct {
	gorm.Model
	Name        string
	Description string
	Price       int
	UmkmID      uint
	IsReady     *bool
	ImgPath     string
}

type MenuParam struct {
	ID     uint   `uri:"menu_id" json:"id"`
	Name   string `form:"name" json:"name" gorm:"-"`
	UmkmID uint   `uri:"umkm_id" form:"umkm_id" json:"umkm_id"`
}

type CreateMenuParam struct {
	Name        string `binding:"required"`
	Description string `binding:"required"`
	Price       int    `binding:"required"`
}

type UpdateMenuParam struct {
	Name        string
	Description string
	Price       int
	IsReady     *bool  `json:"is_ready"`
	ImgPath     string `json:"img_path"`
}
