package models

import (
	"time"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// songに紐づくオーディオファイル
type UserSongAudio struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	UserSongId uint           `gorm:"not null;unique" json:"user_song_id"`
	Name       string         `gorm:"not null" json:"Name"`
	Url        Url            `gorm:"-:all" json:"url"`
	CreatedAt  time.Time      `json:"-"`
	UpdatedAt  time.Time      `json:"-"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// UserSongSection経由で取得、更新するのでメソッド全て不要？
func (audio *UserSongAudio) Create(db *gorm.DB) error {
	result := db.Create(&audio)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (audio *UserSongAudio) GetByID(db *gorm.DB, id uint) error {
	result := db.First(&audio, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (audio *UserSongAudio) GetAllByUserId(db *gorm.DB, userId uint) []UserSongSection {
	var loops []UserSongSection
	result := db.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (audio *UserSongAudio) Update(db *gorm.DB) error {
	result := db.Model(&audio).Debug().Save(audio)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (audio *UserSongAudio) delete(db *gorm.DB) {
	db.Delete(&audio, audio.ID)
}
