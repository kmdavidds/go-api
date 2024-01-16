package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Name        string
	Feedback    string
	Points      uint
	HasReferral bool
	Vouchers    uint
	IsMember    bool
	MemberType  string
	MemberUntil time.Time
	MemberDays  uint
	Address     string
	Kecamatan   string
}
