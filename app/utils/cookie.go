package utils

import (
	"net/http"
	"time"

	"example.com/app/conf"
)

func GetSessionCookie(sessionString string, duration time.Duration) *http.Cookie {
	return &http.Cookie{
		Name:     conf.SESSION_ID_KEY,
		Path:     "/",
		Value:    sessionString,
		Expires:  time.Now().Add(duration),
		MaxAge:   int(duration.Seconds()),
		SameSite: http.SameSiteNoneMode,
		//Domain:   conf.COOKIE_DOMAIN, //環境変数から読み込む
		HttpOnly: true,
		Secure:   true,
	}
}
