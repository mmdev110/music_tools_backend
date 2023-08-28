package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
)

func (h *HandlersConf) TagHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Method)
	if r.Method == http.MethodPost {
		//新規作成、更新
		h.saveTags(w, r)
	} else if r.Method == http.MethodGet {
		//タグ一覧の取得
		h.getTags(w, r)
	} else {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed", r.Method))
	}

}

/*
現在の、タグを全て送ってそれに合わせてテーブルを更新する(宣言的？な)仕様は
タグが増加するにつれてトラフィックが増加していくのでよろしくないかも

/tags/list
/tags/append
/tags/delete
/tags/update
みたいに処理毎にエンドポイントを細分化した方が良さげ

ただ1回のリクエストで完結するため、トランザクションが保たれる点はメリットでもある
インプットを下の様に分けるのが良いか？？

	{
		add:[]Tags
		delete:[]Tags
		update:[]Tags
	}
*/
func (h *HandlersConf) saveTags(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r.Context())
	var input = []models.UserTag{}

	json.NewDecoder(r.Body).Decode(&input)
	//for _, v := range input {
	//	utils.PrintStruct(v)
	//}
	if len(input) == 0 {
		utils.ErrorJSON(w, customError.Others, errors.New("no tags to save"))
		return
	}
	tmp := models.UserTag{}
	db, err := tmp.GetAllByUserId(h.DB, user.ID)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}

	//DBにあって、リクエストにないタグを削除
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
func (h *HandlersConf) getTags(w http.ResponseWriter, r *http.Request) {
	user := h.getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	fmt.Println("@@@gettags")

	//DBから取得
	var tag = models.UserTag{}
	tags, err := tag.GetAllByUserId(h.DB, user.ID)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
	}

	utils.ResponseJSON(w, tags, http.StatusOK)
}
