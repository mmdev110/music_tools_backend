package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/app/models"
	"example.com/app/utils"
)

func TagHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		//新規作成、更新
		saveTags(w, r)
	} else if r.Method == http.MethodGet {
		//タグ一覧の取得
		getTags(w, r)
	} else {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed", r.Method))
	}

}

func saveTags(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	fmt.Println("@@@savetags")
	var input = []models.UserTag{}

	json.NewDecoder(r.Body).Decode(&input)
	//for _, v := range input {
	//	utils.PrintStruct(v)
	//}
	tmp := models.UserTag{}
	db, err := tmp.GetAllByUserId(user.ID)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//タグリレーション削除
	removedTags := utils.FindRemoved(db, input)
	for _, t := range removedTags {
		if err := t.Delete(); err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
	//タグ追加、更新
	models.DB.Debug().Save(&input)
	utils.ResponseJSON(w, input, http.StatusOK)
}
func getTags(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	fmt.Println("@@@gettags")

	//DBから取得
	var tag = models.UserTag{}
	tags, err := tag.GetAllByUserId(user.ID)
	if err != nil {
		utils.ErrorJSON(w, err)
	}

	utils.ResponseJSON(w, tags, http.StatusOK)
}
