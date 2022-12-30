package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"example.com/app/conf"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	fmt.Println(r.Cookies())
	cookie := &http.Cookie{
		Name: "originalcookie",
		//Path:    "/",
		Value:    "lalala",
		Expires:  time.Now().Add(conf.REFRESH_DURATION),
		MaxAge:   int(conf.REFRESH_DURATION.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost", //環境変数から読み込む
		HttpOnly: true,
		Secure:   true,
	}
	http.SetCookie(w, cookie)
	response := map[string]string{"Status": "test"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	js, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		fmt.Println(err)
	}
	w.Write(js)
}
