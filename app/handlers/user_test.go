package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/models"
	"example.com/app/testutil"
)

func Test_UserHandler(t *testing.T) {
	h.DB = TestDB.Begin()
	u, err := models.InsertTestUsersOnly(h.DB)
	if err != nil {
		t.Error(err)
	}
	defer h.DB.Rollback()

	tests := []struct {
		name string
		user *models.User
		code int
	}{
		{"can get requested user", u[0], http.StatusOK},
		//{"cannot get other user", u[1].ID, http.StatusBadRequest},
		{"cannot get unregistered user", &models.User{ID: uint(1)}, http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			r := httptest.NewRequest(http.MethodGet, ts.URL+"/user", nil)
			token, _ := test.user.FakeGenerateToken()
			testutil.AddAuthorizationHeader(r, token)
			r.RequestURI = ""

			w, err2 := ts.Client().Do(r)
			if err2 != nil {
				t.Error(err2)
			}
			defer w.Body.Close()

			//ステータスチェック
			testutil.Checker(t, "status_code", w.StatusCode, test.code)
		})
	}

}
