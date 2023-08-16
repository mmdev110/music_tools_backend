package models

import (
	"time"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// userSongSectionに紐づいたmidiファイル
type UserSectionMidi struct {
	ID                uint   `gorm:"primarykey" json:"id"`
	UserSongSectionId uint   `gorm:"not null" json:"user_song_section_id"`
	Name              string `gorm:"not null" json:"Name"`
	Url               Url    `gorm:"-:all" json:"url"`
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRootsCSV string         `json:"midi_roots_csv"`
	CreatedAt    time.Time      `json:"-"`
	UpdatedAt    time.Time      `json:"-"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// GET用のURLとPUT用のURL
type Url struct {
	Get string `json:"get"`
	Put string `json:"put"`
}

func (midi *UserSectionMidi) Create(db *gorm.DB) error {
	result := db.Create(&midi)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (midi *UserSectionMidi) GetByID(db *gorm.DB, id uint) error {
	result := db.First(&midi, id)
	if result.RowsAffected == 0 {
		return nil
	}
	return nil
}
func (midi *UserSectionMidi) GetAllByUserId(db *gorm.DB, userId uint) []UserSongSection {
	var loops []UserSongSection
	result := db.Where("user_id = ?", userId).Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	return loops
}
func (midi *UserSectionMidi) Update(db *gorm.DB) error {
	result := db.Model(&midi).Updates(midi)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (midi *UserSectionMidi) delete(db *gorm.DB) {
	db.Delete(&midi, midi.ID)
}

func (midi *UserSectionMidi) SetPlaylistName(str string) {
	midi.Name = str
}
