package models

import (
	"time"

	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 曲の各セクション
// イントロ、Aメロなど
type UserAudioRange struct {
	ID                uint           `gorm:"primarykey" json:"id"`
	UserSongSectionId uint           `gorm:"not null" json:"user_song_section_id"`
	Name              string         `json:"name"`
	Start             int            `gorm:"not null" json:"start"`
	End               int            `gorm:"not null" json:"end"`
	SortOrder         int            `gorm:"not null;default:0" json:"sort_order"`
	CreatedAt         time.Time      `json:"-"`
	UpdatedAt         time.Time      `json:"-"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (r UserAudioRange) GetID() uint {
	return r.ID
}

func (r *UserAudioRange) Delete() error {
	//tag自体の削除
	result := DB.Debug().Delete(r)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
