package models

import (
	"errors"

	"gorm.io/gorm"
)

type Book struct {
	gorm.Model
	Title       string  `json:"title" validate:"required,min=3,max=32"`
	Price       float32 `json:"price" validate:"required,gt=0"`
	ImgUrl      string  `json:"imgUrl" validate:"required,min=10"`
	Description string  `json:"description"`
	UserID      uint    `json:"userId" validate:"required"`
	Author      User    `json:"author" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"-"`
}

func (book *Book) BeforeCreate(db *gorm.DB) error {
	user := new(User)
	db.First(user, book.UserID)
	if user.Email == "" {
		return errors.New("invalid user")
	}
	return nil
}

func (book *Book) Create(db *gorm.DB) error {
	if result := db.Create(book); result.Error != nil {
		return result.Error
	}
	return nil
}
