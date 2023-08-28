package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	ID       string `gorm:"primaryKey;type:string;default:uuid_generate_v4()"`
	Email    string `gorm:"size:255;not null;unique"`
	Password string `gorm:"size:255;"`
	Name     string `gorm:"size:255;"`
	Author   []Book `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Carts    []Book `gorm:"many2many:user_carts;joinForeignKey:UserID;joinReferences:BookID"`  //books in carts
	Owns     []Book `gorm:"many2many:user_owns;joinForeignKey:UserID;joinReferences:BookID"`   //books in owns
	Rating   []Book `gorm:"many2many:user_rating;joinForeignKey:UserID;joinReferences:BookID"` //books rated
	Coins    float32
	Role     string
}

// Define the User foreign key relationship for Owns
type UserOwns struct {
	UserID string `gorm:"primaryKey;autoIncrement:false"`
	BookID string `gorm:"primaryKey;autoIncrement:false"`
}

// Define the User foreign key relationship for Cart
type UserCarts struct {
	UserID string `gorm:"primaryKey;autoIncrement:false"`
	BookID string `gorm:"primaryKey;autoIncrement:false"`
}

// Define the User foreign key relationship for book rating
type UserRating struct {
	UserID string `gorm:"primaryKey;autoIncrement:false"`
	BookID string `gorm:"primaryKey;autoIncrement:false"`
}
