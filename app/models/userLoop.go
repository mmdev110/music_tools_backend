package models

import (
	"encoding/json"
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

// clientとの通信に使うstruct
type UserLoopInput struct {
	ID            uint          `json:"id"`
	Name          string        `json:"name"`
	Progressions  []string      `json:"progressions"`
	Key           int           `json:"key"`
	Scale         string        `json:"scale"`
	Memo          string        ` json:"memo"`
	MidiRoots     []int         `json:"midi_roots"`
	UserLoopAudio UserLoopAudio `json:"user_loop_audio"`
	UserLoopMidi  UserLoopMidi  `json:"user_loop_midi"`
	UserLoopTags  []UserLoopTag `json:"user_loop_tags"`
}

// DBに格納するためのstruct
// UserLoopInputの配列要素をstring化している
type UserLoop struct {
	ID     uint `gorm:"primarykey"`
	UserId uint `gorm:"not null"`
	Name   string
	//コード進行をcsv化したもの
	//["Am7","","","Dm7"]->"Am7,,,Dm7"
	Progressions string
	Key          int
	Scale        string
	//オーディオファイル
	UserLoopAudio UserLoopAudio
	//midiファイル
	UserLoopMidi UserLoopMidi
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRoots    string
	Memo         string
	UserLoopTags []UserLoopTag `gorm:"many2many:userloops_tags"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (ul *UserLoop) Create() error {
	result := DB.Omit("UserLoopTags.*").Create(&ul)
	if result.Error != nil {
		return result.Error
	}
	ul.SetMediaUrl()
	return nil
}
func (ul *UserLoop) GetByID(id uint) *gorm.DB {
	result := DB.Model(&UserLoop{}).Preload("UserLoopAudio").Preload("UserLoopMidi").Preload("UserLoopTags").Debug().First(&ul, id)
	if result.RowsAffected == 0 {
		return result
	}
	err := ul.SetMediaUrl()
	if err != nil {
		fmt.Println(err)
	}
	return result
}

type ULSearchCond struct {
	TagIds    []uint `json:"tag_ids"`
	SubString string `json:"substring"`
}

func (ul *UserLoop) GetByUserId(userId uint, condition ULSearchCond) ([]UserLoop, error) {
	var loops []UserLoop
	var result *gorm.DB
	//うまいやり方を考える
	if len(condition.TagIds) > 0 {
		result = DB.Debug().Model(&UserLoop{}).Preload("UserLoopAudio").Preload("UserLoopMidi").Preload("UserLoopTags").Joins("INNER JOIN userloops_tags ult ON user_loops.id=ult.user_loop_id").Joins("INNER JOIN user_loop_tags tags ON tags.id=ult.user_loop_tag_id").Where("user_loops.user_id=? AND tags.id IN ?", userId, condition.TagIds).Find(&loops)
	} else {
		result = DB.Model(&UserLoop{}).Preload("UserLoopAudio").Preload("UserLoopMidi").Preload("UserLoopTags").Debug().Where("user_id=?", userId).Find(&loops)
	}
	if result.RowsAffected == 0 {
		return nil, nil
	}
	if result.Error != nil {
		return nil, result.Error
	}
	for i := range loops {
		err := loops[i].SetMediaUrl()
		if err != nil {
			fmt.Println(err)
		}
	}
	return loops, nil
}
func (ul *UserLoop) Update() error {
	//result := DB.Model(&ul).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Updates(ul)
	result := DB.Debug().Omit("UserLoopTags.*").Save(&ul)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (ul *UserLoop) Delete() error {
	//relationの削除
	err := ul.DeleteTagRelations(ul.UserLoopTags)
	if err != nil {
		return err
	}
	//audio,midiもまとめて削除
	result := DB.Debug().Delete(&ul, ul.ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// input->DB
func (ul *UserLoop) ApplyULInputToUL(ulInput UserLoopInput) {
	prog, _ := json.Marshal(ulInput.Progressions)
	midiroots, _ := json.Marshal(ulInput.MidiRoots)
	//ul.ID = ulInput.ID
	ul.Name = ulInput.Name
	ul.Progressions = string(prog)
	ul.Key = ulInput.Key
	ul.Scale = ulInput.Scale
	ul.MidiRoots = string(midiroots)
	ul.Memo = ulInput.Memo
	ul.UserLoopAudio.Name = ulInput.UserLoopAudio.Name
	ul.UserLoopMidi.Name = ulInput.UserLoopMidi.Name
	ul.UserLoopTags = ulInput.UserLoopTags
	ul.SetMediaUrl()
}

// DB->input
func (uli *UserLoopInput) ApplyULtoULInput(ul UserLoop) {
	var prog []string
	json.Unmarshal([]byte(ul.Progressions), &prog)
	var midiroots []int
	json.Unmarshal([]byte(ul.MidiRoots), &midiroots)
	uli.ID = ul.ID
	uli.Name = ul.Name
	uli.Progressions = prog
	uli.Key = ul.Key
	uli.Scale = ul.Scale
	uli.MidiRoots = midiroots
	uli.Memo = ul.Memo
	uli.UserLoopAudio = ul.UserLoopAudio
	uli.UserLoopMidi = ul.UserLoopMidi
	uli.UserLoopTags = ul.UserLoopTags
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (ul *UserLoop) SetMediaUrl() error {
	Backend := conf.BACKEND_URL

	//audio
	if ul.UserLoopAudio.Name != "" {
		audio := &ul.UserLoopAudio
		//拡張子削除
		//n := strings.ReplaceAll(audio.Name, ".wav", "")
		//n = strings.ReplaceAll(audio.Name, ".mp3", "")
		//hlsName := n + PlaylistSuffix + ".m3u8"
		folder := strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/"
		//hlsPath := folder + hlsName
		//https://{CloudFront_Domain}/{user_id}/{userLoop_id}/{Name}_hls.m3u8
		//ul.UserLoopAudio.Url = backend + "/" + strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + name + PlaylistSuffix + ".m3u8"
		//{backend_domain}/hls/{user_loop_id}
		//m3u8ファイルをバックエンドで書き換える必要があるためGETのurlはバックエンドを指定する
		//get, _ := utils.GenerateSignedUrl(hlsPath, http.MethodGet, conf.PRESIGNED_DURATION)
		//audio.Url.Get = backend + "/" + "hls" + "/" + strconv.Itoa(int(ul.ID))
		audio.Url.Get = Backend + "/hls/" + strconv.Itoa(int(ul.ID))
		put, err := utils.GenerateSignedUrl(folder+audio.Name, http.MethodPut, conf.PRESIGNED_DURATION)
		if err != nil {
			return err
		}
		audio.Url.Put = put
	} else {
		ul.UserLoopAudio.Url.Get = ""
		ul.UserLoopAudio.Url.Put = ""
	}
	//midi
	//midiはpresigned urlを直接返す
	if ul.UserLoopMidi.Name != "" {
		midi := &ul.UserLoopMidi
		path := strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + midi.Name
		//https://{CloudFront_Domain}/{user_id}/{userLoop_id}/{midi.name}
		//ul.UserLoopMidi.Url = "https://" + CFDomain + "/" + strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + midi.Name
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
	} else {
		ul.UserLoopMidi.Url.Get = ""
		ul.UserLoopMidi.Url.Put = ""
	}
	return nil
}

func (ul *UserLoop) GetHLSName() string {
	audio := &ul.UserLoopAudio
	n := strings.ReplaceAll(audio.Name, ".wav", "")
	n = strings.ReplaceAll(audio.Name, ".mp3", "")

	return n + PlaylistSuffix + ".m3u8"
}
func (ul *UserLoop) GetFolderName() string {
	folder := strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/"
	return folder
}

// 中間テーブルのrelationを削除
func (ul *UserLoop) DeleteTagRelations(tags []UserLoopTag) error {
	if len(tags) == 0 {
		return nil
	}
	//中間テーブルのレコード削除
	err := DB.Debug().Model(&ul).Association("UserLoopTags").Delete(tags)
	if err != nil {
		return err
	}
	return nil
}
