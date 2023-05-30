package models

import (
	"time"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userSongSectionに紐づいたmidiファイル
type UserSectionMidi struct {
	ID                uint   `gorm:"primarykey" json:"id"`
	UserSongSectionId uint   `gorm:"not null" json:"user_loop_id"`
	Name              string `gorm:"not null" json:"Name"`
	Url               Url    `gorm:"-:all" json:"url"`
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRootsCSV string `json:"midi_roots_csv"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

// GET用のURLとPUT用のURL
type Url struct {
	Get string `json:"get"`
	Put string `json:"put"`
}

func (midi *UserSectionMidi) Create() error {
	result := DB.Create(&midi)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (midi *UserSectionMidi) GetByID(id uint) error {
	result := DB.First(&midi, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (midi *UserSectionMidi) GetAllByUserId(userId uint) []UserSongSection {
	var loops []UserSongSection
	result := DB.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (midi *UserSectionMidi) Update() error {
	result := DB.Model(&midi).Debug().Updates(midi)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (midi *UserSectionMidi) delete() {
	DB.Delete(&midi, midi.ID)
}

func (midi *UserSectionMidi) SetPlaylistName() {

	midi.Name = midi.Name
}
