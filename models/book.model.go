package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	UserID      uint
	Name        string `gorm:"size:255;"`
	Description string `gorm:"size:255;"`
	Tag         string
	Price       float32
}
