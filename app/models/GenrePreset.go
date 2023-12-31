package models

import (
	"time"

	"gorm.io/gorm"
)

// 固定で用意しているジャンル名
// ユーザーは自分でgenreを入力するか、ここから文字列をコピーして使用することができる
type GenrePreset struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ほぼ更新しないデータなのでとりあえず直書きで定義しておく
var GenrePresets = []string{
	"Chill Wave",
	"Synth Wave",
	"J-POP",
	"Rock",
	"RnB",
	"Funk",
}
