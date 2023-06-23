package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// userSongの一覧
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
		return
	}
	fmt.Println("listhandler")
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)

	//検索条件取り出し
	var condition = models.SongSearchCond{}
	json.NewDecoder(r.Body).Decode(&condition)
	utils.PrintStruct(condition)

	var us = models.UserSong{}
	userSongs, _ := us.GetByUserId(user.ID, condition)
	fmt.Println("list handler response")
	utils.PrintStruct(userSongs)
	utils.ResponseJSON(w, userSongs, http.StatusOK)
}

func SongHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("@@@@songhandler")
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	param := strings.TrimPrefix(r.URL.Path, "/song/")
	fmt.Printf("param = %s\n", param)

	if r.Method == http.MethodPost {
		if param == "new" {
			//新規作成
			createSong(w, r, user)
		} else {
			//更新
			uid, _ := strconv.Atoi(param)
			userSongId := uint(uid)
			updateSong(w, r, user, userSongId)
		}
	} else if r.Method == http.MethodGet {
		//取得
		//uid, _ := strconv.Atoi(param)
		//userSongId := uint(uid)
		uuid := param
		getSong(w, r, user, uuid)
	} else {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
		return
	}

}
func createSong(w http.ResponseWriter, r *http.Request, user *models.User) {
	fmt.Println("@@@@Create Song")
	var us = models.UserSong{}
	json.NewDecoder(r.Body).Decode(&us)
	utils.PrintStruct(us)

	//create
	us.UserId = user.ID
	//uuid付与
	us.UUID = uuid.NewString()
	for {
		err := us.Create()
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			//uuidが衝突してるので更新して再実行
			us.UUID = uuid.NewString()
			continue
		}
		if err != nil {
			utils.ErrorJSON(w, err)
			break
		}
		break
	}

	fmt.Println("@@@@CreateSong response")
	utils.PrintStruct(us)
	utils.ResponseJSON(w, us, http.StatusOK)

}
func updateSong(w http.ResponseWriter, r *http.Request, user *models.User, userSongId uint) {
	fmt.Println("@@@@Update Song")

	var us = models.UserSong{}
	json.NewDecoder(r.Body).Decode(&us)
	utils.PrintStruct(us)

	//update
	var db = models.UserSong{}
	result := db.GetByID(userSongId)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Song not found"))
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
	}

	//タグの中間テーブルの削除
	removedTags := utils.FindRemoved(db.Tags, us.Tags)
	for _, tag := range removedTags {
		if err := db.DeleteTagRelation(&tag); err != nil {
			utils.ErrorJSON(w, err)
		}
	}
	//ジャンルの中間テーブルの削除
	removedGenres := utils.FindRemoved(db.Genres, us.Genres)
	for _, genre := range removedGenres {
		if err := db.DeleteGenreRelation(&genre); err != nil {
			utils.ErrorJSON(w, err)
		}
	}
	//instrumentsの削除
	removedInstruments := utils.FindRemoved(db.Instruments, us.Instruments)
	fmt.Println("removed Instruments: ", len(removedInstruments))
	for _, inst := range removedInstruments {
		if err := inst.Delete(); err != nil {
			utils.ErrorJSON(w, err)
		}
	}
	//audioRangeの削除
	for _, sec := range us.Sections {
		for _, secDB := range db.Sections {
			if sec.ID == secDB.ID {
				removedRange := utils.FindRemoved(secDB.AudioRanges, sec.AudioRanges)
				for _, r := range removedRange {
					if err := r.Delete(); err != nil {
						utils.ErrorJSON(w, err)
					}
				}
			}
		}
	}
	//sectionsの削除
	removedSections := utils.FindRemoved(db.Sections, us.Sections)
	for _, sec := range removedSections {
		if err := sec.Delete(); err != nil {
			utils.ErrorJSON(w, err)
		}
	}
	//section-instrumentsの中間テーブルの削除
	for _, sec := range us.Sections {
		for _, secDB := range db.Sections {
			if sec.ID == secDB.ID {
				removedInst := utils.FindRemoved(secDB.Instruments, sec.Instruments)
				for _, inst := range removedInst {
					sec.DeleteInstrumentRelation(&inst)
				}
			}
		}
	}
	if err := us.Update(); err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//再取得
	//us.GetByID(userSongId)

	//presignedURLセット
	if err := us.SetMediaUrls(); err != nil {
		utils.ErrorJSON(w, err)
	}

	fmt.Println("@@@@UpdateSong response")
	utils.PrintStruct(us)
	utils.ResponseJSON(w, us, http.StatusOK)
}

// songIdに対応するsongを返す
func getSong(w http.ResponseWriter, r *http.Request, user *models.User, uuid string) {

	//DBから取得
	var us = models.UserSong{}
	//result := us.GetByID(userSongId)
	result := us.GetByUUID(uuid)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Song not found"))
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
	}
	//他人のデータは取得不可
	if us.UserId != user.ID {
		utils.ErrorJSON(w, errors.New("you cannot get this Song"))
	}
	utils.ResponseJSON(w, us, http.StatusOK)
}
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
		return
	}
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	type Req struct {
		ID uint `json:"Id"`
	}

	var req = Req{}
	json.NewDecoder(r.Body).Decode(&req)

	us := &models.UserSong{}
	result := us.GetByID(req.ID)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("Song not found"))
		return
	}
	err := us.Delete()
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	type Res struct {
		Message string `json:"message"`
	}

	utils.ResponseJSON(w, &Res{"OK"}, http.StatusOK)
}
