package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	UserID      uint
	Name        string `gorm:"size:255;"`
	Description string `gorm:"size:255;"`
	Tag         string
	Owners      []*User `gorm:"many2many:user_books"`
	InCart      []*User `gorm:"many2many:user_books_cart"`
	Price       float32
}
