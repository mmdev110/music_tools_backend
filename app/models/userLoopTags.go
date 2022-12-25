package models

import (
	"time"

	"example.com/app/utils"
	"gorm.io/gorm"
)

type UserLoopTag struct {
	ID        uint       `gorm:"primarykey" json:"id"`
	UserId    uint       `gorm:"index:idx_uid_name,unique;not null" json:"user_id"`
	Name      string     `gorm:"index:idx_uid_name,unique;not null" json:"name"`
	SortOrder int        `gorm:"not null;default:0" json:"sort_order"`
	UserLoops []UserLoop `gorm:"many2many:userloops_tags" json:"user_loops"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (tag *UserLoopTag) Create() error {
	result := DB.Create(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserLoopTag) Update() error {
	result := DB.Save(&tag)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (tag *UserLoopTag) GetById(id uint) error {
	result := DB.Debug().Preload("UserLoops").First(tag, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// UserLoopsも取得する版
func (tag *UserLoopTag) GetAllByUserId(uid uint) ([]UserLoopTag, error) {
	var uls []UserLoopTag
	result := DB.Debug().Preload("UserLoops").Where("user_id=?", uid).Find(&uls)
	if result.Error != nil {
		return nil, result.Error
	}
	return uls, nil
}

// tagと中間テーブルのrelationを削除
func (tag *UserLoopTag) DeleteTagAndRelations(tags []UserLoopTag) error {
	if len(tags) == 0 {
		return nil
	}
	utils.PrintStruct(tags)
	//中間テーブルのレコード削除
	//err := DB.Debug().Model(&UserLoopTag{}).Association("UserLoops").Delete(tags)
	//if err != nil {
	//	return err
	//}
	//tagの削除
	result := DB.Debug().Model(&UserLoopTag{}).Delete(tags)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
