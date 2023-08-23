package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/models"
	"example.com/app/utils"
)

func Test_UserHandler(t *testing.T) {
	h.DB = TestDB.Begin()
	defer h.DB.Rollback()
	users, err := models.PrepareTestUsersOnly(h.DB, false)
	if err != nil {
		t.Error(err)
	}
	user := users[0]
	r := httptest.NewRequest(http.MethodGet, "/user", nil)
	ctx := utils.SetUIDInContext(r.Context(), user.ID)

	w := httptest.NewRecorder()

	handler := http.HandlerFunc(h.UserHandler)
	handler.ServeHTTP(w, r.WithContext(ctx))

	fmt.Println(utils.BodyToString(w.Result().Body))
}
