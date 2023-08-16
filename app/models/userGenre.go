package models

import (
	"time"

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

func (genre *UserGenre) Create(db *gorm.DB) error {
	result := db.Create(genre)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserGenre) Update(db *gorm.DB) error {
	result := db.Save(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserGenre) GetById(db *gorm.DB, id uint) error {
	result := db.Preload("UserSongs").First(tag, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UserSongsも取得する版
func (tag *UserGenre) GetAllByUserId(db *gorm.DB, uid uint) ([]UserGenre, error) {
	var uls []UserGenre
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
func (g *UserGenre) Delete(db *gorm.DB) error {
	//中間テーブルのレコード削除
	err := db.Model(g).Association("UserSongs").Clear()
	if err != nil {
		return err
	}
	//tagの削除
	result := db.Model(&UserGenre{}).Delete(g)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (g UserGenre) GetID() uint {
	return g.ID
}
