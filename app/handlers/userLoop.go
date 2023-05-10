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
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
		return
	}
	fmt.Println("listhandler")
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)

	//検索条件取り出し
	var condition = models.ULSearchCond{}
	json.NewDecoder(r.Body).Decode(&condition)
	utils.PrintStruct(condition)

	var ul = models.UserLoop{}
	userLoops, _ := ul.GetByUserId(user.ID, condition)
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
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
		return
	}

}

func saveLoop(w http.ResponseWriter, r *http.Request) {
	fmt.Println("@@@@saveLoop")
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	param := strings.TrimPrefix(r.URL.Path, "/loop/")
	fmt.Printf("param = %s\n", param)

	var ulInput = models.UserLoopInput{}
	json.NewDecoder(r.Body).Decode(&ulInput)
	utils.PrintStruct(ulInput)
	var ul = models.UserLoop{}
	utils.PrintStruct(ulInput)

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
		uid, _ := strconv.Atoi(param)
		userLoopId := uint(uid)
		result := ul.GetByID(userLoopId)
		if result.RowsAffected == 0 {
			utils.ErrorJSON(w, errors.New("loop not found"))
		}
		if result.Error != nil {
			utils.ErrorJSON(w, result.Error)
		}
		ul_db := ul
		utils.PrintStruct(ulInput)
		ul.ApplyULInputToUL(ulInput)
		utils.PrintStruct(ul)
		//タグの削除
		//DBにあってinputに存在しないものが対象
		tags_to_delete := []models.UserLoopTag{}
		for _, tag_db := range ul_db.UserLoopTags {
			found := false
			for _, tag := range ul.UserLoopTags {
				if tag_db.ID == tag.ID {
					found = true
				}
			}
			if !found {
				tags_to_delete = append(tags_to_delete, tag_db)
			}
		}
		if err := ul.DeleteTagRelations(tags_to_delete); err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err2 := ul.Update()
		if err2 != nil {
			utils.ErrorJSON(w, err2)
			return
		}
		//再取得
		ul.GetByID(userLoopId)
	}

	responseUri := models.UserLoopInput{}
	responseUri.ApplyULtoULInput(ul)

	utils.ResponseJSON(w, responseUri, http.StatusOK)
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
	result := ul.GetByID(userLoopId)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("loop not found"))
		return
	}
	if result.Error != nil {
		utils.ErrorJSON(w, result.Error)
		return
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
func DeleteLoop(w http.ResponseWriter, r *http.Request) {
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

	ul := &models.UserLoop{}
	result := ul.GetByID(req.ID)
	if result.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("loop not found"))
		return
	}
	err := ul.Delete()
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	type Res struct {
		Message string `json:"message"`
	}

	utils.ResponseJSON(w, &Res{"OK"}, http.StatusOK)
}
