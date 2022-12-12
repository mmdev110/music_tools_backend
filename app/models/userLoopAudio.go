package models

import (
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userLoopに登録されたオーディオファイル
type UserLoopAudio struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	UserLoopId uint   `gorm:"not null" json:"user_loop_id"`
	Name       string `gorm:"not null" json:"Name"`
	Url        Url    `gorm:"-:all" json:"url"`
	gorm.Model
}

// UserLoop経由で取得、更新するのでメソッド全て不要？
func (audio *UserLoopAudio) Create() error {
	result := DB.Create(&audio)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (audio *UserLoopAudio) GetByID(id uint) error {
	result := DB.First(&audio, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (audio *UserLoopAudio) GetAllByUserId(userId uint) []UserLoop {
	var loops []UserLoop
	result := DB.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (audio *UserLoopAudio) Update() error {
	result := DB.Model(&audio).Debug().Updates(audio)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (audio *UserLoopAudio) delete() {
	DB.Delete(&audio, audio.ID)
}
