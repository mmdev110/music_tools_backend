package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
	"gorm.io/gorm"
)

// userSongの一覧
func ListHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed", r.Method))
		return
	}
	fmt.Println("listhandler")
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)

	//検索条件取り出し
	var condition = models.SongSearchCond{}
	json.NewDecoder(r.Body).Decode(&condition)

	//自分のuserId以外は検索禁止
	if ids := condition.UserIds; !(len(ids) == 1 && ids[0] == user.ID) {
		utils.ErrorJSON(w, customError.OperationNotAllowed, errors.New("invalid user_id condition"))
		return
	}
	var us = models.UserSong{}
	userSongs, err := us.Search(DB, condition)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//audioのget情報のみ付与
	for _, v := range userSongs {
		if err := v.SetAudioUrlGet(); err != nil {
			utils.ErrorJSON(w, customError.Others, err)
		}
	}
	fmt.Println("list handler response")
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
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed", r.Method))
		return
	}

}
func createSong(w http.ResponseWriter, r *http.Request, user *models.User) {
	fmt.Println("@@@@Create Song")
	var us = models.UserSong{}
	json.NewDecoder(r.Body).Decode(&us)

	//create
	us.UserId = user.ID

	if err := us.Create(DB); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
	}

	//presignedurlセット
	if err := us.SetMediaUrls(); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
	}

	fmt.Println("@@@@CreateSong response")
	utils.ResponseJSON(w, us, http.StatusOK)

}
func updateSong(w http.ResponseWriter, r *http.Request, user *models.User, userSongId uint) {
	fmt.Println("@@@@Update Song")

	var us = models.UserSong{}
	json.NewDecoder(r.Body).Decode(&us)

	//update
	var db = models.UserSong{}

	err := DB.Debug().Transaction(func(tx *gorm.DB) error {
		//for update
		result := db.GetByID(tx, userSongId, true)
		if result.RowsAffected == 0 {
			return errors.New("Song not found")
		}
		if result.Error != nil {
			return result.Error
		}
		//タグの中間テーブルの削除
		removedTags := utils.FindRemoved(db.Tags, us.Tags)
		for _, tag := range removedTags {
			if err := db.DeleteTagRelation(tx, &tag); err != nil {
				return err
			}
		}
		//ジャンルの中間テーブルの削除
		removedGenres := utils.FindRemoved(db.Genres, us.Genres)
		for _, genre := range removedGenres {
			if err := db.DeleteGenreRelation(tx, &genre); err != nil {
				return err
			}
		}
		//instrumentsの削除
		removedInstruments := utils.FindRemoved(db.Instruments, us.Instruments)
		fmt.Println("removed Instruments: ", len(removedInstruments))
		for _, inst := range removedInstruments {
			if err := inst.Delete(tx); err != nil {
				return err
			}
		}
		//audioRangeの削除
		for _, sec := range us.Sections {
			for _, secDB := range db.Sections {
				if sec.ID == secDB.ID {
					removedRange := utils.FindRemoved(secDB.AudioRanges, sec.AudioRanges)
					for _, r := range removedRange {
						if err := r.Delete(tx); err != nil {
							return err
						}
					}
				}
			}
		}
		//sectionsの削除
		removedSections := utils.FindRemoved(db.Sections, us.Sections)
		for _, sec := range removedSections {
			if err := sec.Delete(tx); err != nil {
				return err
			}
		}
		//section-instrumentsの中間テーブルの削除
		for _, sec := range us.Sections {
			for _, secDB := range db.Sections {
				if sec.ID == secDB.ID {
					removedInst := utils.FindRemoved(secDB.Instruments, sec.Instruments)
					for _, inst := range removedInst {
						if err := sec.DeleteInstrumentRelation(tx, &inst); err != nil {
							return err
						}
					}
				}
			}
		}
		us.LastModifiedAt = time.Now()
		us.LastViewedAt = db.LastViewedAt
		if err := us.Update(tx); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}

	//再取得
	//us.GetByID(userSongId)

	//presignedURLセット
	if err := us.SetMediaUrls(); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}

	utils.ResponseJSON(w, us, http.StatusOK)
}

// songIdに対応するsongを返す
func getSong(w http.ResponseWriter, r *http.Request, user *models.User, uuid string) {

	//DBから取得
	var us = models.UserSong{}
	//result := us.GetByID(userSongId)
	result := us.GetByUUID(DB, uuid, true)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, customError.Others, errors.New("Song not found"))
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, customError.Others, result.Error)
		return
	}
	//他人のデータは取得不可
	if us.UserId != user.ID {
		utils.ErrorJSON(w, customError.Others, errors.New("you cannot get this Song"))
	}
	//閲覧回数の更新
	us.ViewTimes += 1
	us.LastViewedAt = time.Now()
	if err := us.Update(DB); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//presignedurlセット
	//audioのgetのみ
	if err := us.SetAudioUrlGet(); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
	}

	utils.ResponseJSON(w, us, http.StatusOK)
}
func DeleteSong(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed", r.Method))
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
	result := us.GetByID(DB, req.ID, false)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, customError.Others, errors.New("Song not found"))
		return
	}
	err := us.Delete(DB)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	type Res struct {
		Message string `json:"message"`
	}

	utils.ResponseJSON(w, &Res{"OK"}, http.StatusOK)
}
