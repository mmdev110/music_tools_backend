package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"example.com/app/utils"
	_ "golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// clientとの通信に使うstruct
type UserLoopInput struct {
	ID            uint          `json:"id"`
	Progressions  []string      `json:"progressions"`
	Key           int           `json:"key"`
	Scale         string        `json:"scale"`
	Memo          string        ` json:"memo"`
	MidiRoots     []int         `json:"midi_roots"`
	UserLoopAudio UserLoopAudio `json:"user_loop_audio"`
	UserLoopMidi  UserLoopMidi  `json:"user_loop_midi"`
}

// DBに格納するためのstruct
// UserLoopInputの配列要素をstring化している
type UserLoop struct {
	ID     uint `gorm:"primarykey" json:"id"`
	UserId uint `gorm:"not null" json:"user_id"`
	//コード進行をcsv化したもの
	//["Am7","","","Dm7"]->"Am7,,,Dm7"
	Progressions string `json:"progressions"`
	Key          int    `json:"key"`
	Scale        string `json:"scale"`
	//オーディオファイル
	UserLoopAudio UserLoopAudio `json:"user_loop_audio"`
	//midiファイル
	UserLoopMidi UserLoopMidi `json:"midi_path"`
	//midiファイル内でルートとなるノートのindexをcsv化したもの
	//[1,2,3,4]->"1,2,3,4"
	MidiRoots string ` json:"midi_roots"`
	Memo      string ` json:"memo"`
	gorm.Model
}

func (ul *UserLoop) Create() error {
	result := DB.Create(&ul)
	if result.Error != nil {
		return result.Error
	}
	ul.SetMediaUrl()
	return nil
}
func (ul *UserLoop) GetByID(id uint) error {
	result := DB.Model(&UserLoop{}).Preload("UserLoopAudio").Preload("UserLoopMidi").Debug().First(&ul, id)
	if result.RowsAffected == 0 {
		return nil
	}
	err := ul.SetMediaUrl()
	if err != nil {
		fmt.Println(err)
	}
	return nil
}
func (ul *UserLoop) GetAllByUserId(userId uint) []UserLoop {
	var loops []UserLoop
	result := DB.Model(&UserLoop{}).Preload("UserLoopAudio").Preload("UserLoopMidi").Debug().Find(&loops)
	if result.RowsAffected == 0 {
		return nil
	}
	for _, ul := range loops {
		err := ul.SetMediaUrl()
		if err != nil {
			fmt.Println(err)
		}
	}
	return loops
}
func (ul *UserLoop) Update() error {
	result := DB.Model(&ul).Session(&gorm.Session{FullSaveAssociations: true}).Debug().Updates(ul)
	if err := result.Error; err != nil {
		return err
	}
	return nil
}
func (ul *UserLoop) delete() {
	DB.Delete(&ul, ul.ID)
}

func (ul *UserLoop) ApplyULInputToUL(ulInput UserLoopInput) {
	prog, _ := json.Marshal(ulInput.Progressions)
	midiroots, _ := json.Marshal(ulInput.MidiRoots)
	//ul.ID = ulInput.ID
	ul.Progressions = string(prog)
	ul.Key = ulInput.Key
	ul.Scale = ulInput.Scale
	ul.MidiRoots = string(midiroots)
	ul.Memo = ulInput.Memo
	ul.UserLoopAudio.Name = ulInput.UserLoopAudio.Name
	ul.UserLoopMidi.Name = ulInput.UserLoopMidi.Name
	ul.SetMediaUrl()
}
func (uli *UserLoopInput) ApplyULtoULInput(ul UserLoop) {
	var prog []string
	json.Unmarshal([]byte(ul.Progressions), &prog)
	var midiroots []int
	json.Unmarshal([]byte(ul.MidiRoots), &midiroots)
	uli.ID = ul.ID
	uli.Progressions = prog
	uli.Key = ul.Key
	uli.Scale = ul.Scale
	uli.MidiRoots = midiroots
	uli.Memo = ul.Memo
	uli.UserLoopAudio = ul.UserLoopAudio
	uli.UserLoopMidi = ul.UserLoopMidi
}

var PlaylistSuffix = "_hls"

// s3ファイルの格納場所を返す
func (ul *UserLoop) SetMediaUrl() error {
	//CFDomain := os.Getenv("AWS_LOUDFRONT_DOMAIN")
	Backend := os.Getenv("BACKEND_DOMAIN")

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
		//get, _ := utils.GenerateSignedUrl(hlsPath, http.MethodGet, 60*15)
		//audio.Url.Get = backend + "/" + "hls" + "/" + strconv.Itoa(int(ul.ID))
		audio.Url.Get = Backend + "/hls/" + strconv.Itoa(int(ul.ID))
		put, err := utils.GenerateSignedUrl(folder+audio.Name, http.MethodPut, 60*15)
		if err != nil {
			return err
		}
		audio.Url.Put = put
	}
	//midi
	//midiはpresigned urlを直接返す
	if ul.UserLoopMidi.Name != "" {
		midi := &ul.UserLoopMidi
		path := strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + midi.Name
		//https://{CloudFront_Domain}/{user_id}/{userLoop_id}/{midi.name}
		//ul.UserLoopMidi.Url = "https://" + CFDomain + "/" + strconv.Itoa(int(ul.UserId)) + "/" + strconv.Itoa(int(ul.ID)) + "/" + midi.Name
		get, err := utils.GenerateSignedUrl(path, http.MethodGet, 60*15)
		if err != nil {
			return err
		}
		put, err2 := utils.GenerateSignedUrl(path, http.MethodPut, 60*15)
		if err != nil {
			return err2
		}
		ul.UserLoopMidi.Url.Get = get
		ul.UserLoopMidi.Url.Put = put
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
