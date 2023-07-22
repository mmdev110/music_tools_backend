package models

import (
	"time"

	"gorm.io/gorm"
)

// 固定で用意しているタグ名
// ユーザーは自分でtagを入力するか、ここから文字列をコピーして使用することができる
type TagPreset struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	Name      string         `gorm:"not null" json:"name"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// ほぼ更新しないデータなのでとりあえず直書きで定義しておく
var TagPresets = []string{
	"slow",
	"fast",
}
