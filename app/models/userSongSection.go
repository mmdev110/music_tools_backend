package models

import (
	"time"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 曲の各セクション
// イントロ、Aメロなど
type UserSongSection struct {
	ID         uint `gorm:"primarykey" json:"id"`
	UserSongId uint `gorm:"not null" json:"user_song_id"`
	//セクション名
	Name string `json:"section"`
	//コード進行をcsv化したもの
	//["Am7","","","Dm7"]->"Am7,,,Dm7"
	ProgressionsCSV string `json:"progressions_csv"`
	Key             int    `json:"key"`
	BPM             int    `json:"bpm"`
	Scale           string `json:"scale"`
	Memo            string `json:"memo"`
	//オーディオ再生範囲
	LoopRange `json:"audio_playback_range"`
	//midiファイル
	Midi      UserSectionMidi `json:"midi"`
	SortOrder int             `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type LoopRange struct {
	Start int `gorm:"not null" json:"start"`
	End   int `gorm:"not null" json:"end"`
}

func (sec *UserSongSection) Delete() error {
	result := DB.Debug().Delete(sec)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sec UserSongSection) GetID() uint {
	return sec.ID
}
