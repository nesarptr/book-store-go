package models

import "gorm.io/gorm"

type Book struct {
	gorm.Model
	Title       string  `json:"title" validate:"required,min=3,max=32"`
	Price       float32 `json:"price" validate:"required,gt=0"`
	ImgUrl      string  `json:"imgUrl" validate:"required,min=10"`
	Description string  `json:"description"`
	UserID      uint    `json:"userId" validate:"required"`
	Author      User    `json:"author" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID;references:ID" validate:"dive"`
}

func (book *Book) AfterCreate(db *gorm.DB) error {
	user := new(User)
	db.First(user, book.UserID)
	user.Books = append(user.Books, *book)
	db.Save(user)
	return nil
}
