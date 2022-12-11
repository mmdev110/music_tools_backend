package models

import (
	_ "golang.org/x/crypto/bcrypt"
)

// userLoopに登録されたオーディオファイル
type UserLoopMidi struct {
	ID         uint   `gorm:"primarykey" json:"id"`
	UserLoopId uint   `gorm:"not null" json:"user_loop_id"`
	Name       string `gorm:"not null" json:"Name"`
	Url        Url    `gorm:"-:all" json:"url"`
}

// GET用のURLとPUT用のURL
type Url struct {
	Get string `json:"get"`
	Put string `json:"put"`
}

func (midi *UserLoopMidi) Create() error {
	result := DB.Create(&midi)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (midi *UserLoopMidi) GetByID(id uint) error {
	result := DB.First(&midi, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (midi *UserLoopMidi) GetAllByUserId(userId uint) []UserLoop {
	var loops []UserLoop
	result := DB.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (midi *UserLoopMidi) Update() error {
	result := DB.Model(&midi).Debug().Updates(midi)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (midi *UserLoopMidi) delete() {
	DB.Delete(&midi, midi.ID)
}

func (midi *UserLoopMidi) SetPlaylistName() {

	midi.Name = midi.Name
}
