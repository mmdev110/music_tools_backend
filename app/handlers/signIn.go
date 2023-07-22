package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"example.com/app/conf"
	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
)

func SignInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("signIn")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed for signin", r.Method))
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
		utils.ErrorJSON(w, customError.UserNotFound, fmt.Errorf("user not found for %s", form.Email))
		return
	}
	if !user.IsConfirmed {
		utils.ErrorJSON(w, customError.AddressNotConfirmed, fmt.Errorf("address %s found, but not confirmed yet.", form.Email))
		return
	}
	//check password
	ok := user.ComparePassword(form.Password)
	if !ok {
		utils.ErrorJSON(w, customError.IncorrectPassword, errors.New("password mismatch"))
		return
	}
	//generate jwt
	accessToken, _ := user.GenerateToken("access", conf.TOKEN_DURATION)
	refreshToken, _ := user.GenerateToken("refresh", conf.REFRESH_DURATION)
	user.AccessToken = accessToken
	if err := user.Update(); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}

	//session生成
	session := models.Session{}
	result := session.GetByUserID(user.ID)
	if result.RowsAffected == 0 {
		session.UserId = user.ID
	}
	session.SessionString = uuid.NewString()
	session.RefreshToken = "Bearer " + refreshToken
	if err := session.Update(); err != nil {
		utils.ErrorJSON(w, customError.Others, err)
		return
	}
	//sessionIdをクッキーにセットさせる
	//httponly, secure, samesite
	cookie := utils.GetSessionCookie(session.SessionString, conf.REFRESH_DURATION)
	http.SetCookie(w, cookie)

	type Response = struct {
		User        *models.User `json:"user"`
		AccessToken string       `json:"access_token"`
	}
	utils.ResponseJSON(w, &Response{user, accessToken}, http.StatusOK)
}
