package entity

import "gorm.io/gorm"

type Umkm struct {
	gorm.Model
	Name   string
	Slogan string
}

type UmkmParam struct {
	ID     uint   `json:"id" uri:"umkm_id"`
	Name   string `json:"name" form:"name" gorm:"-"`
	Slogan string `json:"slogan"`
}

type CreateUmkmParam struct {
	Name   string `binding:"required"`
	Slogan string `binding:"required"`
}

type UpdateUmkmParam struct {
	Name   string
	Slogan string
}
