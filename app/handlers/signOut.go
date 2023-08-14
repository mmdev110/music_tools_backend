package handlers

import (
	"fmt"
	"net/http"

	"example.com/app/customError"
	"example.com/app/utils"
)

// sign out
// refresh_tokenを削除して返す
func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		utils.ErrorJSON(w, customError.MethodNotAllowed, fmt.Errorf("method %s not allowed for SignOut", r.Method))
		return
	}
	//無効なcookieで上書きする
	cookie := utils.GetInvalidCookie()
	http.SetCookie(w, cookie)

	response := struct {
		Message string `json:"message"`
	}{}

	//response
	utils.ResponseJSON(w, response, http.StatusOK)
}
