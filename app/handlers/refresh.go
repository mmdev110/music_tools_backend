package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"example.com/app/conf"
	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
)

/*
cookieによるリフレッシュ
*/
func (h *HandlersConf) RefreshHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("refresh")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, customError.Others, fmt.Errorf("method %s not allowed for refresh", r.Method))
		return
	}
	cookie, err := r.Cookie(conf.SESSION_ID_KEY)
	if err != nil {
		utils.ErrorJSON(w, customError.Others, errors.New("session_id not found"))
		return
	}
	sessionId := cookie.Value

	//sessionIdでsession取得
	session := &models.Session{}
	sessionResult := session.GetBySessionID(h.DB, sessionId)
	if sessionResult.RowsAffected == 0 {
		utils.ErrorJSON(w, customError.Others, errors.New("session not found"))
		return
	} else if sessionResult.Error != nil {
		utils.ErrorJSON(w, customError.Others, sessionResult.Error)
		return
	}
	//refreshTokenを検証
	claim, err := utils.Authenticate(session.RefreshToken, "refresh")
	if err != nil {
		utils.ErrorJSON(w, customError.Others, sessionResult.Error)
		return
	}
	if claim.TokenType != "refresh" {
		utils.ErrorJSON(w, customError.Others, errors.New("invalid tokentype"))
		return
	}
	userId := claim.UserId
	user := models.GetUserByID(h.DB, userId)

	//regenerate jwt
	accessToken, _ := user.GenerateToken("access")
	refreshToken, _ := user.GenerateToken("refresh")
	user.AccessToken = accessToken
	user.Update(h.DB)

	//sessionの更新
	//fmt.Println("old_refreshToken= ", session.RefreshToken)
	session.SessionString = uuid.NewString()
	session.RefreshToken = "Bearer " + refreshToken
	//fmt.Println("new_refreshToken= ", session.RefreshToken)

	if err := session.Update(h.DB); err != nil {
		utils.ErrorJSON(w, customError.Others, errors.New("session save failed"))
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
