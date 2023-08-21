package handlers

import (
	"fmt"
	"net/http"

	"example.com/app/utils"
)

func (h *Base) UserHandler(w http.ResponseWriter, r *http.Request) {
	//動作確認用
	//presignedUrl := awsUtil.GenerateSignedUrl()
	user := h.getUserFromContext(r.Context())
	fmt.Printf("userid in handler = %d\n", user.ID)
	utils.ResponseJSON(w, user, http.StatusOK)
}
