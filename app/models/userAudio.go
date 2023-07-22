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
func (audio *UserSongAudio) Create() error {
	result := DB.Create(&audio)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (audio *UserSongAudio) GetByID(id uint) error {
	result := DB.First(&audio, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (audio *UserSongAudio) GetAllByUserId(userId uint) []UserSongSection {
	var loops []UserSongSection
	result := DB.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (audio *UserSongAudio) Update() error {
	result := DB.Model(&audio).Debug().Save(audio)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (audio *UserSongAudio) delete() {
	DB.Delete(&audio, audio.ID)
}
