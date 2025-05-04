package models

import "gorm.io/gorm"

type Person struct {
	gorm.Model
	Name        string `gorm:"index"`
	Surname     string `gorm:"index"`
	Patronymic  string
	Age         int
	Nationality string
	Gender      string
	IsActive    bool `gorm:"default:true"`
}
