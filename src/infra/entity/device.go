package entity

import (
	"gorm.io/gorm"
)

type Device struct {
	gorm.Model
	Name   string
	Code   string `gorm:"unique"`
	Secret string
	Status int `gorm:"type:bit"`
	Owner  User
}
