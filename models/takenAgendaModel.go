package models

import (
	"gorm.io/gorm"
)

type TakenAgenda struct {
	gorm.Model
	AgendaId       uint
	CreatorEmail   string
	TakerEmail     string
	Kecamatan      string
	Address        string
	IsDone         bool
	OrganicKilo    uint
	AnorganicKilo  uint
	MetalKilo      uint
	ElectronicKilo uint
	OtherKilo      uint
}
