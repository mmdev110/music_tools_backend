package utils

import (
	"net/http"
	"time"

	"example.com/app/conf"
)

func GetSessionCookie(sessionString string, duration time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:    conf.SessionID_KEY,
		Path:    "/",
		Value:   sessionString,
		Expires: time.Now().Add(duration),
		MaxAge:  int(duration.Seconds()),
		//SameSite: http.SameSiteStrictMode,
		Domain: "localhost", //環境変数から読み込む
		//HttpOnly: true,
		//Secure: true,
	}
}
