package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/utils"
)

// sign up
// メールアドレス確認の流れ
// SignUpHandlerでフロントエンドのリンク(トークン付)をメール送信
// ユーザーが踏む
// フロントエンドからEmailConfirmationHandlerを叩く
// 確認できたらsigninページに遷移
// 確認できなかったらエラー文表示
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	fmt.Println("signUp")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed for SignUp", r.Method))
		return
	}
	type Form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)
	//find existing user by email
	existingUser := models.GetUserByEmail(form.Email)
	if existingUser != nil && existingUser.IsConfirmed {
		utils.ErrorJSON(w, fmt.Errorf("user already exists for %s", existingUser.Email))
		return
	}
	var user *models.User
	if existingUser != nil && !existingUser.IsConfirmed {
		//update password
		existingUser.SetNewPassword(form.Password)
		if err := existingUser.Update(); err != nil {
			utils.ErrorJSON(w, fmt.Errorf("error while updating existing user: %v", err))
			return
		}
		user = existingUser
	} else {
		//create new user
		newUser, err := models.CreateUser(form.Email, form.Password)
		if err != nil {
			utils.ErrorJSON(w, fmt.Errorf("error while creating new user: %v", err))
			return
		}
		user = &newUser
	}
	//send confirmation email
	//リセット用のトークン生成
	token, err := user.GenerateToken("email_confirm", 30*time.Minute)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	//リセット用のリンク生成
	link, err2 := url.JoinPath(conf.FRONTEND_URL, "email_confirm")
	if err2 != nil {
		utils.ErrorJSON(w, err2)
		return
	}
	link = link + fmt.Sprintf("?token=%s", token)
	//メール送信
	body := "メールアドレス確認用のリンクをお送りいたします。\n" +
		"30分以内に下記のリンクにアクセスしていただくことでメールアドレスの確認が完了いたします。\n" +
		"上記の内容に心当たりがない場合はこのメールを無視してください。\n" +
		link
	utils.SendEmail(user.Email, "Email Confirmation(music_tools)", body)
	//response
	utils.ResponseJSON(w, user, http.StatusOK)
}

// メールアドレス確認ハンドラ
func EmailConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Token string `json:"token"`
	}
	var req Request
	json.NewDecoder(r.Body).Decode(&req)
	//tokenのチェック
	claims, err := utils.ParseJwt(req.Token)
	if err != nil {
		utils.ErrorJSON(w, err)
		return
	}
	//user取得
	user := models.GetUserByID(claims.UserId)
	if user == nil {
		utils.ErrorJSON(w, errors.New("user not found"))
		return
	}
	if user.IsConfirmed {
		utils.ErrorJSON(w, errors.New("address is already confirmed"))
		return
	}
	//確認完了
	user.IsConfirmed = true
	user.Update()

	utils.ResponseJSON(w, user, http.StatusOK)
}
