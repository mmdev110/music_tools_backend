package models

import (
	"fmt"
	"time"

	"example.com/app/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID          uint   `gorm:"primarykey" json:"user_id"`
	Email       string `gorm:"unique;not null" json:"email"`
	IsConfirmed bool   `gorm:"not null;default:0" json:"is_confirmed"`
	//トークン類はユーザーに返さない
	Password     string         `gorm:"not null" json:"-"`
	AccessToken  string         `json:"-"`
	RefreshToken string         `json:"-"`
	Songs        []UserSong     `json:"songs"`
	Tags         []UserTag      `json:"tags"`
	Genres       []UserGenre    `json:"genres"`
	Session      Session        `json:"-"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}
type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func CreateUser(db *gorm.DB, email, password string) (User, error) {
	user := User{
		Email:       email,
		Password:    encrypt(password),
		IsConfirmed: false,
		//Token:    utils.GenerateJwt(email),
	}
	db.Create(&user)
	return user, nil
}
func GetUserByID(db *gorm.DB, id uint) *User {
	var user User
	result := db.First(&user, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}
func GetUserByEmail(db *gorm.DB, email string) *User {
	var user User
	result := db.First(&user, "Email = ?", email)
	if result.RowsAffected == 0 {
		return nil
	}
	return &user
}
func (user *User) Update(db *gorm.DB) error {
	result := db.Model(&user).Debug().Updates(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (user *User) delete(db *gorm.DB) {
	db.Delete(&user, user.ID)
}

func encrypt(password string) string {
	//TODO
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes)
}
func (user *User) SetNewPassword(password string) error {
	user.Password = encrypt(password)
	return nil
}

func (user *User) ComparePassword(input string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input))
	return err == nil
}
func (user *User) GenerateToken(tokenType string, duration time.Duration) (string, error) {
	AllowedTokenTypes := []string{"access", "reset", "refresh", "email_confirm"}
	isAllowed := false
	for _, v := range AllowedTokenTypes {
		if tokenType == v {
			isAllowed = true
			break
		}
	}
	//allowedtokentypesにない場合はエラー
	if !isAllowed {
		return "", fmt.Errorf("token type: %s not allowed", tokenType)
	}

	token, err := utils.GenerateJwt(user.ID, tokenType, duration)
	if err != nil {
		return "", fmt.Errorf("error while generating token: %v", err)
	}
	return token, nil
}
