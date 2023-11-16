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
	Name             string `binding:"required" json:"name"`
	Slogan           string `binding:"required" json:"slogan"`
	OwnerName        string `binding:"required" json:"owner_name"`
	OwnerPhoneNumber string `binding:"required" json:"owner_phone_number"`
}

type UpdateUmkmParam struct {
	Name             string `json:"name"`
	Slogan           string `json:"slogan"`
	Status           string `json:"status"`
	OwnerName        string `json:"owner_name"`
	OwnerPhoneNumber string `json:"owner_phone_number"`
}
