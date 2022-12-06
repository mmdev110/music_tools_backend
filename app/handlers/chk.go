package handlers

import (
	"net/http"

	"example.com/app/utils"
)

func ChkHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	response := map[string]string{"Status": "OK"}
	utils.ResponseJSON(w, response, http.StatusOK)
}
