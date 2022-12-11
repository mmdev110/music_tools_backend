package models

import (
	"errors"
	"os"
	"strconv"
	"strings"

	_ "golang.org/x/crypto/bcrypt"
)

// userLoopに登録されたオーディオファイル
type UserLoopAudio struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	UserLoopId   uint   `gorm:"not null" json:"user_loop_id"`
	Name         string `gorm:"not null" json:"Name"`
	OriginalName string `gorm:"not null" json:"original_name"`
}

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

func (audio *UserLoopAudio) SetPlaylistName() {
	if audio.OriginalName == "" {
		return
	}
	name := strings.ReplaceAll(audio.OriginalName, ".wav", "")
	name = strings.ReplaceAll(audio.OriginalName, ".mp3", "")
	audio.Name = name
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (audio *UserLoopAudio) GetAudioUrl(user *User, ul *UserLoop) (string, error) {
	if audio.Name == "" {
		return "", errors.New("audio Name not set.")
	}
	CFDomain := os.Getenv("AWS_CLOUDFRONT_DOMAIN")
	return "https://" + CFDomain + "/" + strconv.Itoa(int(user.ID)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + audio.Name + PlaylistSuffix + ".m3u8", nil
}
