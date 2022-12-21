package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"example.com/app/models"
	"example.com/app/utils"
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
	token, _ := user.GenerateToken("access", 24*time.Hour)
	user.Token = token
	user.Update()

	utils.ResponseJSON(w, user, http.StatusOK)
}
