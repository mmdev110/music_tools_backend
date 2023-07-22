package models

import (
	"time"

	"example.com/app/utils"
	"gorm.io/gorm"
)

type UserSongInstrument struct {
	ID         uint              `gorm:"primarykey" json:"id"`
	UserSongId uint              `gorm:"index:idx_inst_uid_name,unique;not null" json:"user_song_id"`
	Name       string            `gorm:"index:idx_inst_uid_name,unique;not null" json:"name"`
	SortOrder  int               `gorm:"not null;default:0" json:"sort_order"`
	Sections   []UserSongSection `gorm:"many2many:sections_instruments" json:"song_sections"`
	Category   string            `gorm:"not null" json:"category"`
	Memo       string            `json:"memo"`
	CreatedAt  time.Time         `json:"-"`
	UpdatedAt  time.Time         `json:"-"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"-"`
}

func (inst *UserSongInstrument) Create() error {
	result := DB.Create(&inst)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (inst *UserSongInstrument) Update() error {
	result := DB.Save(&inst)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// instと中間テーブルのrelationを削除
func (inst *UserSongInstrument) Delete(db *gorm.DB) error {
	if db == nil {
		db = DB
	}
	utils.PrintStruct(inst)
	//中間テーブルのレコード削除
	err := DB.Debug().Model(inst).Association("Sections").Clear()
	if err != nil {
		return err
	}
	//inst自体の削除
	result := DB.Debug().Delete(inst)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (inst UserSongInstrument) GetID() uint {
	return inst.ID
}
