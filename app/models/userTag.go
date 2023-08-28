package models

import (
	"time"

	"gorm.io/gorm"
)

type UserTag struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserId    uint           `gorm:"index:idx_tag_uid_name,unique;not null" json:"user_id"`
	Name      string         `gorm:"index:idx_tag_uid_name,unique;not null" json:"name"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	UserSongs []UserSong     `gorm:"many2many:usersongs_tags" json:"user_songs"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (tag *UserTag) Create(db *gorm.DB) error {
	result := db.Create(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserTag) Update(db *gorm.DB) error {
	result := db.Save(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserTag) GetById(db *gorm.DB, id uint) error {
	result := db.Preload("UserSongs").First(tag, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UserSongsも取得する版
func (tag *UserTag) GetAllByUserId(db *gorm.DB, uid uint) ([]UserTag, error) {
	var uls []UserTag
	result := db.
		Preload("UserSongs").
		Where("user_id=?", uid).
		Order("sort_order ASC").
		Find(&uls)
	if result.Error != nil {
		return nil, result.Error
	}
	return uls, nil
}

// tagと中間テーブルのrelationを削除
func (tag *UserTag) Delete(db *gorm.DB) error {
	//中間テーブルのレコード削除
	err := db.Model(tag).Association("UserSongs").Clear()
	if err != nil {
		return err
	}
	//tag自体の削除
	result := db.Delete(tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (tag UserTag) GetID() uint {
	return tag.ID
}
