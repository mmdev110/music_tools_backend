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
)

// userLoopの一覧
func ListHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	var ul = models.UserLoop{}
	userLoops := ul.GetAllByUserId(user.ID)
	var userLoopsInputs []models.UserLoopInput
	for _, v := range userLoops {
		//UserLoopsをUserLoopsINputsに変換
		uli := models.UserLoopInput{}
		uli.ApplyULtoULInput(v)
		userLoopsInputs = append(userLoopsInputs, uli)
	}
	utils.ResponseJSON(w, userLoopsInputs, http.StatusOK)
}

// userLoopの保存
type S3Url struct {
	Mp3  string `json:"mp3"`
	Midi string `json:"midi"`
}
type LoopHandlerResponse struct {
	UserLoopInput models.UserLoopInput `json:"user_loop_input"`
	S3Url         S3Url                `json:"s3url"`
}

func LoopHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//新規作成、更新
		saveLoop(w, r)
	} else if r.Method == http.MethodGet {
		//取得
		getLoop(w, r)
	} else {
		utils.ErrorJSON(w, errors.New("method not allowed"))
	}

}

func saveLoop(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	param := strings.TrimPrefix(r.URL.Path, "/loop/")
	fmt.Printf("param = %s\n", param)

	var ulInput = models.UserLoopInput{}
	json.NewDecoder(r.Body).Decode(&ulInput)
	utils.PrintStruct(ulInput)
	var ul = models.UserLoop{}

	if param == "new" {
		//create
		ul.UserId = user.ID
		ul.ApplyULInputToUL(ulInput)
		err := ul.Create()
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	} else {
		userLoopId, _ := strconv.Atoi(param)
		err := ul.GetByID(uint(userLoopId))
		if err != nil {
			utils.ErrorJSON(w, err)
		}
		//update
		ul.ApplyULInputToUL(ulInput)
		utils.PrintStruct(ul)

		err2 := ul.Update()
		if err2 != nil {
			utils.ErrorJSON(w, err)
		}
	}
	responseUri := models.UserLoopInput{}
	responseUri.ApplyULtoULInput(ul)
	//s3presignedURL生成(PUT用)
	mp3Url := ""
	midiUrl := ""
	fmt.Println("@@@beforegenerate")
	if ul.AudioPath != "" {
		url, err := utils.GenerateSignedUrl(ul.ID, ul.AudioPath, http.MethodPut, 15*60)
		if err != nil {
			utils.ErrorJSON(w, fmt.Errorf("error at GenerateSignedUrl: %v", err))
			return
		}
		mp3Url = url
	}
	if ul.MidiPath != "" {
		url, err := utils.GenerateSignedUrl(ul.ID, ul.AudioPath, http.MethodPut, 15*60)
		if err != nil {
			utils.ErrorJSON(w, fmt.Errorf("error at GenerateSignedUrl: %v", err))
			return
		}
		midiUrl = url
	}
	fmt.Printf("AudioPath: %s.mp3Url:  %s\n", ul.AudioPath, mp3Url)
	fmt.Printf("MidiPath: %s.midiUrl:  %s\n", ul.MidiPath, midiUrl)
	response := LoopHandlerResponse{
		UserLoopInput: responseUri,
		S3Url:         S3Url{mp3Url, midiUrl},
	}

	utils.ResponseJSON(w, response, http.StatusOK)
}
func getLoop(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	param := strings.TrimPrefix(r.URL.Path, "/loop/")
	userLoopIdInt, _ := strconv.Atoi(param)
	userLoopId := uint(userLoopIdInt)
	fmt.Printf("param = %s\n", param)

	//DBから取得
	var ul = models.UserLoop{}
	err := ul.GetByID(userLoopId)
	if err != nil {
		utils.ErrorJSON(w, err)
	}
	//他人のデータは取得不可
	if ul.UserId != user.ID {
		utils.ErrorJSON(w, errors.New("you cannot get this loop"))
	}
	//s3presignedURL生成(GET用)
	mp3Url := ""
	midiUrl := ""
	if ul.AudioPath != "" {
		url, err := utils.GenerateSignedUrl(ul.ID, ul.AudioPath, http.MethodGet, 15*60)
		if err == nil {
			mp3Url = url
		}
	}
	if ul.MidiPath != "" {
		url, err := utils.GenerateSignedUrl(ul.ID, ul.AudioPath, http.MethodGet, 15*60)
		if err == nil {
			midiUrl = url
		}

	}
	var ulInput = models.UserLoopInput{}
	//ulをuliに変換
	ulInput.ApplyULtoULInput(ul)
	response := LoopHandlerResponse{
		UserLoopInput: ulInput,
		S3Url:         S3Url{mp3Url, midiUrl},
	}

	utils.ResponseJSON(w, response, http.StatusOK)
}
