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
	*LoopRange `json:"audio_playback_range"`
	//midiファイル
	Midi      UserSectionMidi `json:"user_loop_midi"`
	SortOrder int             `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type LoopRange struct {
	Start uint `gorm:"not null" json:"start"`
	End   uint `gorm:"not null" json:"end"`
}
