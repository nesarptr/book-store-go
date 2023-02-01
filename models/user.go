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
	Password string  `json:"-" validate:"required,min=6"`
	Books    []Book  `json:"-" validate:"dive"`
	Orders   []Order `json:"-" validate:"dive"`
}

func (user *User) BeforeCreate(db *gorm.DB) error {
	foundUser := User{}
	db.Where("email = ?", user.Email).First(&foundUser)
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

func (user *User) Create(db *gorm.DB) error {
	if result := db.Create(user); result.Error != nil {
		return result.Error
	}
	return nil
}
