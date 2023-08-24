package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
)

func (h *HandlersConf) GenreHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		//新規作成、更新
		h.saveGenres(w, r)
	} else if r.Method == http.MethodGet {
		//タグ一覧の取得
		h.getGenres(w, r)
	} else {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed", r.Method))
	}

}

func (h *HandlersConf) saveGenres(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	fmt.Println("@@@saveGenres")
	var input = []models.UserGenre{}

	json.NewDecoder(r.Body).Decode(&input)
	//for _, v := range input {
	//	utils.PrintStruct(v)
	//}
	tmp := models.UserGenre{}
	db, err := tmp.GetAllByUserId(h.DB, user.ID)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}

	//h.DBにあって、リクエストにないタグを削除
	removedTags := utils.FindRemoved(db, input)
	for _, t := range removedTags {
		if err := t.Delete(h.DB); err != nil {
			utils.ErrorJSON(w, customError.Others, err)
			return
		}
	}
	//タグ追加、更新
	h.DB.Save(&input)
	utils.ResponseJSON(w, input, http.StatusOK)
}
func (h *HandlersConf) getGenres(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	fmt.Println("@@@gettags")

	//h.DBから取得
	var tag = models.UserGenre{}
	tags, err := tag.GetAllByUserId(h.DB, user.ID)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
	}

	utils.ResponseJSON(w, tags, http.StatusOK)
}
