package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"example.com/app/conf"
	"example.com/app/customError"
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
func (h *Base) SignUpHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.MethodNotAllowed, fmt.Errorf("method %s not allowed for SignUp", r.Method))
		return
	}
	type Form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)

	if form.Email == "" {
		utils.ErrorJSON(w, customError.InsufficientParameters, errors.New("email not provided"))
		return
	} else if form.Password == "" {
		utils.ErrorJSON(w, customError.InsufficientParameters, errors.New("password not provided"))
		return
	}
	//find existing user by email
	existingUser := models.GetUserByEmail(h.DB, form.Email)
	if existingUser != nil && existingUser.IsConfirmed {
		utils.ErrorJSON(w, customError.UserAlreadyExists, nil)
		return
	}
	var user *models.User
	if existingUser != nil && !existingUser.IsConfirmed {
		//update password
		existingUser.SetNewPassword(form.Password)
		if err := existingUser.Update(h.DB); err != nil {
			utils.ErrorJSON(w, customError.Others, fmt.Errorf("error while updating existing user: %v", err))
			return
		}
		user = existingUser
	} else {
		//create new user
		newUser, err := models.CreateUser(h.DB, form.Email, form.Password)
		if err != nil {
			utils.ErrorJSON(w, customError.Others, fmt.Errorf("error while creating new user: %v", err))
			return
		}
		user = &newUser
	}
	//send confirmation email
	//メール確認用のトークン生成
	token, err := user.GenerateToken("email_confirm")
	if err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//メール確認用のリンク生成
	link, err2 := url.JoinPath(conf.FRONTEND_URL, "email_confirm")
	if err2 != nil {
		utils.ErrorJSON(w, customError.Others, err2)
		return
	}
	link = link + fmt.Sprintf("?token=%s", token)
	//メール送信
	body := "メールアドレス確認用のリンクをお送りいたします。\n" +
		"30分以内に下記のリンクにアクセスしていただくことでメールアドレスの確認が完了いたします。\n" +
		"上記の内容に心当たりがない場合はこのメールを無視してください。\n" +
		link
	if h.SendEmail {
		utils.SendEmail(user.Email, "Email Confirmation(music_tools)", body)
	}
	//response
	utils.ResponseJSON(w, user, http.StatusOK)
}

// メールアドレス確認ハンドラ
func (h *Base) EmailConfirmationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.MethodNotAllowed, fmt.Errorf("method %s not allowed for email confirmation", r.Method))
		return
	}
	type Request struct {
		Token string `json:"token"`
	}
	var req Request
	json.NewDecoder(r.Body).Decode(&req)
	if req.Token == "" {
		utils.ErrorJSON(w, customError.InvalidToken, nil)
		return
	}
	//tokenのチェック
	claims, err := utils.ParseJwt(req.Token)
	if err != nil {
		utils.ErrorJSON(w, customError.InvalidToken, err)
		return
	}
	if claims.TokenType != "email_confirm" {
		utils.ErrorJSON(w, customError.InvalidToken, errors.New("wrong token type"))
		return
	}
	//user取得
	user := models.GetUserByID(h.DB, claims.UserId)
	if user == nil {
		utils.ErrorJSON(w, customError.UserNotFound, nil)
		return
	}
	if user.IsConfirmed {
		utils.ErrorJSON(w, customError.AddressAlreadyConfirmed, nil)
		return
	}
	//確認完了
	user.IsConfirmed = true
	user.Update(h.DB)

	utils.ResponseJSON(w, user, http.StatusOK)
}
