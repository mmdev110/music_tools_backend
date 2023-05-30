package models

import (
	"time"

	"gorm.io/gorm"
)

// 固定で用意しているジャンル名
// ユーザーは自分でgenreを入力するか、ここから文字列をコピーして使用することができる
type GenrePreset struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"index:unique;not null" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
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
