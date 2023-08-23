package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255;"`
	Name     string `gorm:"size:255;"`
	Author   []Book //things that the user has posted to sell on the platform
	Owns     []Book `gorm:"many2many:user_books"`
	Cart     []Book `gorm:"many2many:user_books_cart"`
	Coins    float32
	Role     string
}
