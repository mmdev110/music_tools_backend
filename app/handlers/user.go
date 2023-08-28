package handlers

import (
	"net/http"

	"example.com/app/customError"
	"example.com/app/utils"
)

func (h *HandlersConf) UserHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	user := h.getUserFromContext(r.Context())
	if user == nil {
		utils.ErrorJSON(w, customError.UserNotFound, nil)
		return
	}
	//fmt.Printf("userid in handler = %d\n", user.ID)
	utils.ResponseJSON(w, user, http.StatusOK)
}
