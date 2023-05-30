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
	tags, err := tmp.GetAllByUserId(user.ID)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}

	//削除する(=DBにあってinputに無い)レコードを探す
	tags_to_delete := []models.UserTag{}
	tags_not_deleted := []models.UserTag{}
	for _, tag_db := range tags {
		found_in_input := false
		for _, tag_input := range input {
			if tag_db.ID == tag_input.ID {
				found_in_input = true
			}
		}
		if !found_in_input {
			fmt.Println("delete tag")
			utils.PrintStruct(tag_db)
			tags_to_delete = append(tags_to_delete, tag_db)
		} else {
			tags_not_deleted = append(tags_not_deleted, tag_db)
		}
	}
	var t = models.UserTag{}
	if err := t.DeleteTagAndRelations(tags_to_delete); err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	tags = tags_not_deleted

	//更新と追加
	for _, tag_input := range input {
		found_in_db := false
		for _, tag_db := range tags {
			//更新
			if tag_db.ID == tag_input.ID {
				tag_db.Name = tag_input.Name
				fmt.Println("existing tag")
				utils.PrintStruct(tag_db)
				found_in_db = true
			}
		}
		//更新されない(=DB内に無い)tag_inputを追加する
		if !found_in_db {
			tag_input.UserId = user.ID
			fmt.Println("new tag")
			utils.PrintStruct(tag_input)
			tags = append(tags, tag_input)
		}
	}
	models.DB.Debug().Save(tags)
	//再取得して返す
	responseTags, err2 := tmp.GetAllByUserId(user.ID)
	if err2 != nil {
		utils.ErrorJSON(w, err2)
		return
	}
	utils.ResponseJSON(w, responseTags, http.StatusOK)
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
