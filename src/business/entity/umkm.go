package entity

import "gorm.io/gorm"

const (
	StatusOpen     = "open"
	StatusClose    = "close"
	StatusInactive = "inactive"
)

type Umkm struct {
	gorm.Model
	Name             string
	Slogan           string
	Status           string
	OwnerName        string
	OwnerPhoneNumber string
}

type UmkmParam struct {
	ID     uint   `json:"id" uri:"umkm_id"`
	Name   string `json:"name" form:"name" gorm:"-"`
	Slogan string `json:"slogan"`
	Status string `json:"status" form:"status"`
}

type CreateUmkmParam struct {
	Name             string `binding:"required"`
	Slogan           string `binding:"required"`
	OwnerName        string `binding:"required"`
	OwnerPhoneNumber string `binding:"required"`
}

type UpdateUmkmParam struct {
	Name             string
	Slogan           string
	Status           string
	OwnerName        string
	OwnerPhoneNumber string
}
