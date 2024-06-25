package entity

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        string `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Username  string         `gorm:"unique"`
	Password  string
	Devices   []Device
}

func (u *User) HasDevice(deviceCode string) bool {
	for _, device := range u.Devices {
		if device.Code == deviceCode {
			return true
		}
	}
	return false
}
