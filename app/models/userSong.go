package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/app/conf"
	"example.com/app/utils"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// song情報
type UserSong struct {
	ID       uint              `gorm:"primarykey" json:"id"`
	UserId   uint              `gorm:"not null" json:"user_id"`
	Title    string            `json:"title"`
	Artist   string            `json:"artist"`
	Sections []UserSongSection `json:"user_song_section"`
	Memo     string            `json:"memo"`
	//オーディオファイル
	Audio UserSongAudio `json:"user_audio"`
	//ジャンル
	Genres []UserGenre `gorm:"many2many:usersongs_genres" json:"genres"`
	//タグ
	Tags      []UserTag `gorm:"many2many:usersongs_tags" json:"tags"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (us *UserSong) Create() error {
	result := DB.Omit("UserTags.*").Create(&us)
	if result.Error != nil {
		return result.Error
	}
	us.SetMediaUrls()
	return nil
}

// songを返す
func (us *UserSong) GetByID(id uint) *gorm.DB {
	result := DB.Model(&UserSongSection{}).Preload("UserAudio").Preload("UserSongMidi").Preload("UserTags").Debug().First(&us, id)
	if result.RowsAffected == 0 {
		return result
	}
	err := us.SetMediaUrls()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

type ULSearchCond struct {
	TagIds    []uint `json:"tag_ids"`
	SubString string `json:"substring"`
}

// userIdに紐づくsong(検索条件があればそれも考慮する)
func (us *UserSong) GetByUserId(userId uint, condition ULSearchCond) ([]UserSong, error) {
	var songs []UserSong
	var result *gorm.DB
	//うまいやり方を考える
	if len(condition.TagIds) > 0 {
		result = DB.Debug().Model(&UserSongSection{}).Preload("UserAudio").Preload("UserSongMidi").Preload("UserTags").Joins("INNER JOIN userloops_tags ult ON user_loops.id=ult.user_loop_id").Joins("INNER JOIN user_loop_tags tags ON tags.id=ult.user_loop_tag_id").Where("user_loops.user_id=? AND tags.id IN ?", userId, condition.TagIds).Find(&songs)
	} else {
		result = DB.Model(&UserSongSection{}).Preload("UserAudio").Preload("UserSongMidi").Preload("UserTags").Debug().Where("user_id=?", userId).Find(&songs)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	for i := range songs {
		err := songs[i].SetMediaUrls()
		if err != nil {
			fmt.Println(err)
		}
	}
	return songs, nil
}
func (us *UserSong) Update() error {
	fmt.Println("@@@@update")
	//result := DB.Model(&us).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Updates(ul)
	result := DB.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Omit("UserTags.*").Save(&us)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (us *UserSong) Delete() error {
	//relationの削除
	err := us.DeleteTagRelations(us.Tags)
	if err != nil {
		return err
	}
	//audio,midiもまとめて削除
	result := DB.Debug().Delete(&us, us.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (us *UserSong) SetMediaUrls() error {

	//audio
	if err := us.setAudioUrl(); err != nil {
		return err
	}
	//midi
	if err := us.setMidiUrl(); err != nil {
		return err
	}
	return nil
}
func (us *UserSong) setAudioUrl() error {
	Backend := conf.BACKEND_URL
	fmt.Println("@@@@setmediaurl")
	fmt.Println(us.Audio)
	//audio
	if us.Audio.Name != "" {
		audio := &us.Audio

		//put urlはpresigned URL
		//get urlはm3u8ファイルを書き換える必要があるためバックエンドを指定する
		audio.Url.Get = Backend + "/hls/" + strconv.Itoa(int(us.ID))
		fmt.Println(conf.PRESIGNED_DURATION)
		put, err := utils.GenerateSignedUrl(us.GetFolderName()+audio.Name, http.MethodPut, conf.PRESIGNED_DURATION)
		if err != nil {
			return err
		}
		audio.Url.Put = put
	}
	return nil
}

// midiファイルの格納場所
// presigned urlを返す
func (us *UserSong) setMidiUrl() error {
	for _, section := range us.Sections {
		if section.Midi.Name != "" {
			midi := &section.Midi
			path := strconv.Itoa(int(us.UserId)) + "/" + strconv.Itoa(int(section.ID)) + "/" + midi.Name
			get, err := utils.GenerateSignedUrl(path, http.MethodGet, conf.PRESIGNED_DURATION)
			if err != nil {
				return err
			}
			put, err2 := utils.GenerateSignedUrl(path, http.MethodPut, conf.PRESIGNED_DURATION)
			if err != nil {
				return err2
			}
			midi.Url.Get = get
			midi.Url.Put = put
		}
	}
	return nil
}

// m3u8ファイルの名前を返す
// (オーディオファイル)_hls.m3u8というルールになっている
func (us *UserSong) GetHLSName() string {
	audio := &us.Audio
	n := strings.ReplaceAll(audio.Name, ".wav", "")
	n = strings.ReplaceAll(n, ".mp3", "")
	n = strings.ReplaceAll(n, ".m4a", "")
	fmt.Println(n)
	return n + PlaylistSuffix + ".m3u8"
}
func (us *UserSong) GetFolderName() string {
	folder := strconv.Itoa(int(us.UserId)) + "/"
	return folder
}

// 中間テーブルのrelationを削除
func (us *UserSong) DeleteTagRelations(tags []UserTag) error {
	if len(tags) == 0 {
		return nil
	}
	//中間テーブルのレコード削除
	err := DB.Debug().Model(&us).Association("UserTags").Delete(tags)
	if err != nil {
		return err
	}
	return nil
}
