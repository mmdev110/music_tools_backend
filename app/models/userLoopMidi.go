package models

import (
	"errors"
	"os"
	"strconv"

	_ "golang.org/x/crypto/bcrypt"
)

// userLoopに登録されたオーディオファイル
type UserLoopMidi struct {
	ID           uint   `gorm:"primarykey" json:"id"`
	UserLoopId   uint   `gorm:"not null" json:"user_loop_id"`
	Name         string `gorm:"not null" json:"Name"`
	OriginalName string `gorm:"not null" json:"original_name"`
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRoots string ` json:"midi_roots"`
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

// s3ファイルの格納場所を返す
func (midi *UserLoopMidi) GetmidiUrl(user *User, ul *UserLoop) (string, error) {
	if midi.Name == "" {
		return "", errors.New("midi Name not set.")
	}
	CFDomain := os.Getenv("AWS_CLOUDFRONT_DOMAIN")
	return "https://" + CFDomain + "/" + strconv.Itoa(int(user.ID)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + midi.Name, nil
}
