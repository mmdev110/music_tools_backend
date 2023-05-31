package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID            uint   `gorm:"primarykey"`
	SessionString string `gorm:"unique;not null"`
	UserId        uint   `gorm:"unique;not null"`
	RefreshToken  string `gorm:"not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (s *Session) Create(refreshToken string) (*Session, error) {
	result := DB.Create(s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}
func (s *Session) GetByUserID(uid uint) *gorm.DB {
	result := DB.Model(&Session{}).Where("user_id=?", uid).First(s)
	return result
}
func (s *Session) GetBySessionID(sessionString string) *gorm.DB {
	result := DB.Model(&Session{}).Where("session_string=?", sessionString).First(s)
	return result
}
func (s *Session) Update() error {
	result := DB.Save(&s)
	return result.Error
}
