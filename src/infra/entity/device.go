package entity

import (
	"gorm.io/gorm"
	"time"
)

type Device struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string
	Code      string `gorm:"unique"`
	Secret    string
	Status    int `gorm:"type:bit"`
	UserID    uint
}
