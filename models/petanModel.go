package models

import "gorm.io/gorm"

type Petan struct {
	gorm.Model
	UserEmail string
	Name      string
	Kecamatan string
	Address   string
	IsTruk    bool
	IsDone    bool
}
