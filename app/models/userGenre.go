package models

import (
	"time"

	"example.com/app/utils"
	"gorm.io/gorm"
)

type UserGenre struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserId    uint           `gorm:"index:idx_genre_uid_name,unique;not null" json:"user_id"`
	Name      string         `gorm:"index:idx_genre_uid_name,unique;not null" json:"name"`
	SortOrder int            `gorm:"not null;default:0" json:"sort_order"`
	UserSongs []UserSong     `gorm:"many2many:usersongs_genres" json:"user_songs"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (genre *UserGenre) Create() error {
	result := DB.Create(genre)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserGenre) Update() error {
	result := DB.Save(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserGenre) GetById(id uint) error {
	result := DB.Debug().Preload("UserSongs").First(tag, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UserSongsも取得する版
func (tag *UserGenre) GetAllByUserId(uid uint) ([]UserGenre, error) {
	var uls []UserGenre
	result := DB.Debug().
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
func (g *UserGenre) Delete() error {
	utils.PrintStruct(g)
	//中間テーブルのレコード削除
	err := DB.Debug().Model(g).Association("UserSongs").Clear()
	if err != nil {
		return err
	}
	//tagの削除
	result := DB.Debug().Model(&UserGenre{}).Delete(g)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (g UserGenre) GetID() uint {
	return g.ID
}
