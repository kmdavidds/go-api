package models

import "gorm.io/gorm"

type Scavenger struct {
	gorm.Model
	Email       string `gorm:"unique"`
	Password    string
	Name        string
	Points      uint
}