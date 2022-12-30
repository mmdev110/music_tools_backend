package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signIn")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed for signin", r.Method))
		return
	}
	type Form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)

	//getUserByEmail
	user := models.GetUserByEmail(form.Email)
	if user == nil {
		utils.ErrorJSON(w, fmt.Errorf("user not found for %s", form.Email))
		return
	}
	//check password
	ok := user.ComparePassword(form.Password)
	if !ok {
		utils.ErrorJSON(w, errors.New("password mismatch"))
		return
	}
	//generate jwt
	accessToken, _ := user.GenerateToken("access", conf.TOKEN_DURATION)
	refreshToken, _ := user.GenerateToken("refresh", conf.REFRESH_DURATION)
	user.AccessToken = accessToken
	user.Update()

	//session生成
	session := models.Session{}
	result := session.GetByUserID(user.ID)
	if result.RowsAffected == 0 {
		session.UserId = user.ID
	}
	session.SessionString = uuid.NewString()
	session.RefreshToken = "Bearer " + refreshToken
	session.Update()
	//sessionIdをクッキーにセットさせる
	//httponly, secure, samesite
	cookie := utils.GetSessionCookie(session.SessionString, conf.REFRESH_DURATION)
	http.SetCookie(w, cookie)
	fmt.Println("header:")
	fmt.Println(w.Header())

	type Response = struct {
		User        *models.User `json:"user"`
		AccessToken string       `json:"access_token"`
	}
	utils.ResponseJSON(w, &Response{user, accessToken}, http.StatusOK)
}
