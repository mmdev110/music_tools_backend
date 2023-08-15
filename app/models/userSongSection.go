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
	ProgressionsCSV string               `json:"progressions_csv"`
	Instruments     []UserSongInstrument `gorm:"many2many:sections_instruments" json:"instruments"`
	Key             int                  `json:"key"`
	BPM             int                  `json:"bpm"`
	BarLength       int                  `json:"bar_length"`
	Scale           string               `json:"scale"`
	Memo            string               `json:"memo"`
	MemoTransition  string               `json:"memo_transition"`
	//オーディオ再生範囲
	AudioRanges []UserAudioRange `json:"audio_ranges"`
	//midiファイル
	Midi      UserSectionMidi `json:"midi"`
	SortOrder int             `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt gorm.DeletedAt  `gorm:"index" json:"-"`
}

type LoopRange struct {
	Start int `gorm:"not null" json:"start"`
	End   int `gorm:"not null" json:"end"`
}

func (sec *UserSongSection) Create(db *gorm.DB) error {
	//Instrumentsはrelationのみ作成
	result := db.Omit("Instruments.*").Create(&sec)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sec *UserSongSection) Delete(db *gorm.DB) error {
	//中間テーブルのレコード削除
	if err := db.Model(sec).Association("Instruments").Clear(); err != nil {
		return err
	}
	result := db.Delete(sec)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (sec UserSongSection) GetID() uint {
	return sec.ID
}

// 中間テーブルのrelationを削除
func (sec *UserSongSection) DeleteInstrumentRelation(db *gorm.DB, inst *UserSongInstrument) error {
	if err := db.Model(sec).Association("Instruments").Delete(inst); err != nil {
		return err
	}
	return nil
}
