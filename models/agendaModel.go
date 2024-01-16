package models

import (
	"time"

	"gorm.io/gorm"
)

type Agenda struct {
	gorm.Model
	CreatorEmail string
	Kecamatan    string
	Date         time.Time
}
