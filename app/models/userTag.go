package models

import (
	"time"

	"example.com/app/utils"
	"gorm.io/gorm"
)

type UserTag struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserId    uint       `gorm:"index:idx_uid_name,unique;not null" json:"user_id"`
	Name      string     `gorm:"index:idx_uid_name,unique;not null" json:"name"`
	SortOrder int        `gorm:"not null;default:0" json:"sort_order"`
	UserSongs []UserSong `gorm:"many2many:usersongs_tags" json:"user_loops"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (tag *UserTag) Create() error {
	result := DB.Create(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserTag) Update() error {
	result := DB.Save(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserTag) GetById(id uint) error {
	result := DB.Debug().Preload("UserSongSections").First(tag, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UserSongsも取得する版
func (tag *UserTag) GetAllByUserId(uid uint) ([]UserTag, error) {
	var uls []UserTag
	result := DB.Debug().Preload("UserSongs").Where("user_id=?", uid).Find(&uls)
	if result.Error != nil {
		return nil, result.Error
	}
	return uls, nil
}

// tagと中間テーブルのrelationを削除
func (tag *UserTag) Delete() error {
	utils.PrintStruct(tag)
	//中間テーブルのレコード削除
	err := DB.Debug().Model(tag).Association("UserSongs").Clear()
	if err != nil {
		return err
	}
	//tag自体の削除
	result := DB.Debug().Delete(tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (tag UserTag) PrintID() uint {
	return tag.ID
}
