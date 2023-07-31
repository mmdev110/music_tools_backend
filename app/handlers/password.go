package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"example.com/app/conf"
	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
)

// パスワードリセット用のリンクをメールで送信するハンドラー
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	action := r.URL.Query().Get("action")
	if action == "request" {
		SendResetEmailHandler(w, r)
	} else if action == "reset" {
		PasswordResetHandler(w, r)
	} else {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("action %s not allowed for this operation", action))
	}

}
func SendResetEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		if r.Method != http.MethodPost {
			utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed for this operation", r.Method))
			return
		}
	}
	fmt.Println("@@SendResetEmailHandler")
	type Form struct {
		Email string `json:"email"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)
	//emailを持ったユーザーがいるか確認
	user := models.GetUserByEmail(DB, form.Email)
	if user == nil {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("user not found for %s", form.Email))
		return
	}
	//リセット用のトークン生成
	token, err := user.GenerateToken("reset", 30*time.Minute)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//リセット用のリンク生成
	link, err2 := url.JoinPath(conf.FRONTEND_URL, "reset_password", "new")
	if err2 != nil {
		utils.ErrorJSON(w, customError.Others, err2)
		return
	}
	link = link + fmt.Sprintf("?token=%s", token)
	//メール送信
	body := "パスワードリセット用のリンクをお送りいたします。\n" +
		"30分以内に下記のリンクより新しいパスワードを設定してください。\n" +
		"上記の内容に心当たりがない場合はこのメールを無視してください。\n" +
		link
	utils.SendEmail(user.Email, "Password Reset(music_tools)", body)
	//
	type Response struct {
		Message string `json:"message"`
	}
	utils.ResponseJSON(w, Response{Message: "OK"}, http.StatusOK)
}

// 新しいパスワードを設定するハンドラー
func PasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	//form取り出し
	type Form struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)
	//tokenのチェック
	claims, err := utils.ParseJwt(form.Token)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//user取得
	user := models.GetUserByID(DB, claims.UserId)
	if user == nil {
		utils.ErrorJSON(w, customError.Others, errors.New("user not found"))
		return
	}
	//新しいパスワード生成
	err3 := user.SetNewPassword(form.NewPassword)
	if err3 != nil {
		utils.ErrorJSON(w, customError.Others, err3)
		return
	}
	user.Update(DB)
	//user更新

	utils.ResponseJSON(w, user, http.StatusOK)
}
