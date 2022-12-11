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

type LoopHandlerResponse struct {
	UserLoopInput models.UserLoopInput `json:"user_loop_input"`
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
		//update
		userLoopId, _ := strconv.Atoi(param)
		err := ul.GetByID(uint(userLoopId))
		if err != nil {
			utils.ErrorJSON(w, err)
		}
		ul.ApplyULInputToUL(ulInput)
		//utils.PrintStruct(ul)

		err2 := ul.Update()
		if err2 != nil {
			utils.ErrorJSON(w, err)
		}
	}
	responseUri := models.UserLoopInput{}
	responseUri.ApplyULtoULInput(ul)

	response := LoopHandlerResponse{
		UserLoopInput: responseUri,
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
	var ulInput = models.UserLoopInput{}
	//ulをuliに変換
	ulInput.ApplyULtoULInput(ul)
	response := LoopHandlerResponse{
		UserLoopInput: ulInput,
	}

	utils.ResponseJSON(w, response, http.StatusOK)
}
