package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	NamaPengguna string `json:"nama_pengguna" validate:"required,min=3"`
	KataSandi    string `json:"kata_sandi" validate:"required,min=6"`
	Role         string `json:"role" validate:"required"`
	RefreshToken string `json:"refresh_token"`
}
