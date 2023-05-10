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
//
//	type UserLoopInput struct {
//		ID            uint              `json:"id"`
//		Name          string            `json:"name"`
//		Progressions  []string          `json:"progressions"`
//		Key           int               `json:"key"`
//		Scale         string            `json:"scale"`
//		Memo          string            ` json:"memo"`
//		UserLoopAudio UserLoopAudio     `json:"user_loop_audio"`
//		UserLoopMidi  UserLoopMidiInput `json:"user_loop_midi"`
//		UserLoopTags  []UserLoopTag     `json:"user_loop_tags"`
//		UserLoop
//	}
type UserLoopInput struct {
	UserLoop
	Progressions []string          `json:"progressions"`
	UserLoopMidi UserLoopMidiInput `json:"user_loop_midi"`
}

// DBに格納するためのstruct
// UserLoopInputの配列要素をstring化している
type UserLoop struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	UserId  uint   `gorm:"not null" json:"user_id"`
	Name    string `json:"name"`
	Artist  string `json:"artist"`
	Section string `json:"section"`
	//コード進行をcsv化したもの
	//["Am7","","","Dm7"]->"Am7,,,Dm7"
	Progressions   string
	Key            int    `json:"key"`
	BPM            int    `json:"bpm"`
	Scale          string `json:"scale"`
	Memo           string `json:"memo"`
	MemoBass       string `json:"memo_nass"`
	MemoChord      string `json:"memo_chord"`
	MemoLead       string `json:"memo_lead"`
	MemoRhythm     string `json:"memo_rhythm"`
	MemoTransition string `json:"memo_transition"`
	//オーディオファイル
	UserLoopAudio UserLoopAudio
	//midiファイル
	UserLoopMidi UserLoopMidi
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
	fmt.Println("@@@@update")
	//result := DB.Model(&ul).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Updates(ul)
	result := DB.Debug().Session(&gorm.Session{FullSaveAssociations: true}).Omit("UserLoopTags.*").Save(&ul)
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
	midi := ulInput.UserLoopMidi
	prog, _ := json.Marshal(ulInput.Progressions)
	midiroots, _ := json.Marshal(midi.MidiRoots)
	//ul.ID = ulInput.ID
	ul.Name = ulInput.Name
	ul.Progressions = string(prog)
	ul.Key = ulInput.Key
	ul.BPM = ulInput.BPM
	ul.Section = ulInput.Section
	ul.Scale = ulInput.Scale
	ul.Memo = ulInput.Memo
	ul.MemoBass = ulInput.MemoBass
	ul.MemoChord = ulInput.MemoChord
	ul.MemoLead = ulInput.MemoLead
	ul.MemoRhythm = ulInput.MemoRhythm
	ul.MemoTransition = ulInput.MemoTransition
	ul.UserLoopAudio = ulInput.UserLoopAudio
	ul.UserLoopMidi = UserLoopMidi{
		ID:   midi.ID,
		Name: midi.Name,
		Url: Url{
			Get: midi.Url.Get,
			Put: midi.Url.Put,
		},
		MidiRoots: string(midiroots),
	}
	ul.UserLoopTags = ulInput.UserLoopTags
	ul.SetMediaUrl()
}

// DB->input
func (uli *UserLoopInput) ApplyULtoULInput(ul UserLoop) {
	var prog []string
	midi := ul.UserLoopMidi
	json.Unmarshal([]byte(ul.Progressions), &prog)
	var midiroots []int
	json.Unmarshal([]byte(midi.MidiRoots), &midiroots)
	uli.ID = ul.ID
	uli.Name = ul.Name
	uli.Progressions = prog
	uli.Key = ul.Key
	uli.BPM = ul.BPM
	uli.Section = ul.Section
	uli.Scale = ul.Scale
	uli.Memo = ul.Memo
	uli.MemoBass = ul.MemoBass
	uli.MemoLead = ul.MemoLead
	uli.MemoChord = ul.MemoChord
	uli.MemoRhythm = ul.MemoRhythm
	uli.MemoTransition = ul.MemoTransition
	uli.UserLoopAudio = ul.UserLoopAudio
	uli.UserLoopMidi = UserLoopMidiInput{
		UserLoopMidi: UserLoopMidi{
			ID:   midi.ID,
			Name: midi.Name,
			Url: Url{
				Get: midi.Url.Get,
				Put: midi.Url.Put,
			},
		},
		MidiRoots: midiroots,
	}
	uli.UserLoopTags = ul.UserLoopTags
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (ul *UserLoop) SetMediaUrl() error {
	Backend := conf.BACKEND_URL
	fmt.Println("@@@@setmediaurl")
	fmt.Println(ul.UserLoopAudio)
	//audio
	if ul.UserLoopAudio.Name != "" {
		audio := &ul.UserLoopAudio
		//拡張子削除
		//n := strings.ReplaceAll(audio.Name, ".wav", "")
		//n = strings.ReplaceAll(audio.Name, ".mp3", "")
		//hlsName := n + PlaylistSuffix + ".m3u8"
		//hlsPath := folder + hlsName
		//https://{CloudFront_Domain}/{user_id}/{userLoop_id}/{Name}_hls.m3u8
		//ul.UserLoopAudio.Url = backend + "/" + strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + name + PlaylistSuffix + ".m3u8"
		//{backend_domain}/hls/{user_loop_id}
		//m3u8ファイルをバックエンドで書き換える必要があるためGETのurlはバックエンドを指定する
		//get, _ := utils.GenerateSignedUrl(hlsPath, http.MethodGet, conf.PRESIGNED_DURATION)
		//audio.Url.Get = backend + "/" + "hls" + "/" + strconv.Itoa(int(ul.ID))
		audio.Url.Get = Backend + "/hls/" + strconv.Itoa(int(ul.ID))
		fmt.Println(conf.PRESIGNED_DURATION)
		put, err := utils.GenerateSignedUrl(ul.GetFolderName()+audio.Name, http.MethodPut, conf.PRESIGNED_DURATION)
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
	n = strings.ReplaceAll(n, ".mp3", "")
	n = strings.ReplaceAll(n, ".m4a", "")
	fmt.Println(n)
	return n + PlaylistSuffix + ".m3u8"
}
func (ul *UserLoop) GetFolderName() string {
	folder := strconv.Itoa(int(ul.UserId)) + "/"
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
