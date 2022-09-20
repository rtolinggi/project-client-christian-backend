package models

import (
	"gorm.io/gorm"
)

type Karyawan struct {
	gorm.Model
	NamaLengkap   string `json:"nama_lengkap" validate:"required"`
	NomorRekening string `json:"nomor_rekening" validate:"required"`
	NamaRekening  string `json:"nama_rekening" validate:"required"`
	NamaBank      string `json:"nama_bank" validate:"required"`
	BankCabang    string `json:"bank_cabang" validtae:"required"`
	Status_aktif  bool   `json:"status_aktif" validate:"required"`
	UserID        uint
	User          User `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
