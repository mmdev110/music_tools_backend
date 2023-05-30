package models

import (
	"time"

	"gorm.io/gorm"
)

// 固定で用意しているタグ名
// ユーザーは自分でtagを入力するか、ここから文字列をコピーして使用することができる
type TagPreset struct {
	ID        uint   `gorm:"primarykey" json:"id"`
	Name      string `gorm:"index:unique;not null" json:"name"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ほぼ更新しないデータなのでとりあえず直書きで定義しておく
var TagPresets = []string{
	"slow",
	"fast",
}
