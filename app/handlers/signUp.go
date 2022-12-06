package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"example.com/app/models"
	"example.com/app/utils"
)

func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	fmt.Println("signUp")
	if r.Method != http.MethodPost {
		utils.ErrorJSON(w, fmt.Errorf("method %s not allowed for SignUp", r.Method))
		return
	}
	type Form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var form Form
	json.NewDecoder(r.Body).Decode(&form)
	//find existing user by email
	existingUser := models.GetUserByEmail(form.Email)
	if existingUser != nil {
		utils.ErrorJSON(w, fmt.Errorf("user already exists for %s", existingUser.Email))
		return
	}
	//create new user
	newUser, err := models.CreateUser(form.Email, form.Password)
	if err != nil {
		utils.ErrorJSON(w, fmt.Errorf("error while creating new user: %v", err))
		return
	}
	//response
	utils.ResponseJSON(w, newUser, http.StatusOK)
}
