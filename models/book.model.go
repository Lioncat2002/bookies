package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	ID          string `gorm:"primaryKey;type:string;default:uuid_generate_v4()"`
	UserID      string
	Name        string `gorm:"size:255;"`
	Description string `gorm:"size:255;"`
	Tag         string
	Price       float32
}
