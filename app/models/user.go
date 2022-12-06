package models

import (
	"fmt"

	"example.com/app/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uint   `gorm:"primarykey" json:"user_id"`
	Email     string `gorm:"unique;not null" json:"email"`
	Password  string `gorm:"not null" json:"-"`
	Token     string `gorm:"not null" json:"token"`
	UserLoops []UserLoop
	gorm.Model
}

func CreateUser(email, password string) (User, error) {
	user := User{
		Email:    email,
		Password: encrypt(password),
		//Token:    utils.GenerateJwt(email),
	}
	DB.Create(&user)
	return user, nil
}
func GetUserByID(id uint) *User {
	var user User
	result := DB.First(&user, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}
func GetUserByEmail(email string) *User {
	var user User
	result := DB.First(&user, "Email = ?", email)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}
func (user *User) update() {
	DB.Model(&user).Debug().Updates(user)
}
func (user *User) delete() {
	DB.Delete(&user, user.ID)
}

func encrypt(password string) string {
	//TODO
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes)
}

func (user *User) ComparePassword(input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input))
	return err == nil
}
func (user *User) GenerateToken() (string, error) {
	token, err := utils.GenerateJwt(user.ID)
	if err != nil {
		return "", fmt.Errorf("error while generating token: %v", err)
	}
	//user.update(struct{ Token string }{token})
	user.Token = token
	user.update()
	return token, nil
}
