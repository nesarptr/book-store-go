package models

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string  `json:"name" validate:"required,min=3,max=32"`
	Email    string  `json:"email" validate:"required,email"`
	Password string  `json:"password" validate:"required,min=6"`
	Books    []Book  `json:"books" validate:"dive"`
	Orders   []Order `json:"orders" validate:"dive"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	foundUser := User{}
	tx.Where("email = ?", user.Email).First(&foundUser)
	if foundUser.Name != "" {
		return errors.New("User Already Exist with this email")
	}
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPw)
	return nil
}
