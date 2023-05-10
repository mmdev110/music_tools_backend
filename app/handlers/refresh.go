package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
)

func RefreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("refresh")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed for refresh", r.Method))
		return
	}
	cookie, err := r.Cookie(conf.SESSION_ID_KEY)
	if err != nil {
		utils.ErrorJSON(w, errors.New("session_id not found"))
		return
	}
	sessionId := cookie.Value

	//sessionIdでsession取得
	session := &models.Session{}
	sessionResult := session.GetBySessionID(sessionId)
	if sessionResult.RowsAffected == 0 {
		utils.ErrorJSON(w, errors.New("session not found"))
		return
	} else if sessionResult.Error != nil {
		utils.ErrorJSON(w, sessionResult.Error)
		return
	}
	//refreshTokenを検証
	claim, err := utils.Authenticate(session.RefreshToken, "refresh")
	if err != nil {
		utils.ErrorJSON(w, sessionResult.Error)
		return
	}
	if claim.TokenType != "refresh" {
		utils.ErrorJSON(w, errors.New("invalid tokentype"))
		return
	}
	userId := claim.UserId
	user := models.GetUserByID(userId)

	//regenerate jwt
	accessToken, _ := user.GenerateToken("access", conf.TOKEN_DURATION)
	refreshToken, _ := user.GenerateToken("refresh", conf.REFRESH_DURATION)
	user.AccessToken = accessToken
	user.Update()

	//sessionの更新
	fmt.Println("old_refreshToken= ", session.RefreshToken)
	session.SessionString = uuid.NewString()
	session.RefreshToken = "Bearer " + refreshToken
	fmt.Println("new_refreshToken= ", session.RefreshToken)

	if err := session.Update(); err != nil {
		utils.ErrorJSON(w, errors.New("session save failed"))
		return
	}

	//sessionIdをクッキーにセットさせる
	//httponly, secure, samesite
	newCookie := utils.GetSessionCookie(session.SessionString, conf.REFRESH_DURATION)
	http.SetCookie(w, newCookie)

	type Response = struct {
		AccessToken string `json:"access_token"`
	}
	utils.ResponseJSON(w, &Response{accessToken}, http.StatusOK)
}