package models

import (
	"time"

	"gorm.io/gorm"
)

type Session struct {
	ID            uint           `gorm:"primarykey"`
	SessionString string         `gorm:"unique;not null"`
	UserId        uint           `gorm:"unique;not null"`
	RefreshToken  string         `gorm:"not null"`
	CreatedAt     time.Time      `json:"-"`
	UpdatedAt     time.Time      `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (s *Session) Create(db *gorm.DB, refreshToken string) (*Session, error) {
	result := db.Create(s)
	if result.Error != nil {
		return nil, result.Error
	}
	return s, nil
}
func (s *Session) GetByUserID(db *gorm.DB, uid uint) *gorm.DB {
	result := db.Model(&Session{}).Where("user_id=?", uid).First(s)
	return result
}
func (s *Session) GetBySessionID(db *gorm.DB, sessionString string) *gorm.DB {
	result := db.Model(&Session{}).Where("session_string=?", sessionString).First(s)
	return result
}
func (s *Session) Update(db *gorm.DB) error {
	result := db.Save(&s)
	return result.Error
}
