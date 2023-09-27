package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/models"
	"example.com/app/testutil"
	"example.com/app/utils"
)

func Test_SaveGenres(t *testing.T) {

}
func Test_GetGenres(t *testing.T) {

	t.Run("can get genres", func(t *testing.T) {
		h.DB = TestDB.Begin()
		defer h.DB.Rollback()

		data, err := models.InsertTestData(h.DB)
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest(http.MethodGet, ts.URL+"/genres", nil)
		req.RequestURI = ""
		token, _ := data.User.FakeGenerateToken()
		testutil.AddAuthorizationHeader(req, token)

		res, err := ts.Client().Do(req)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		testutil.Checker(t, "status_code", res.StatusCode, http.StatusOK)
		var res_genres []models.UserGenre
		if err := utils.BodyToStruct(res.Body, &res_genres); err != nil {
			t.Error(err)
		}

		testutil.Checker(t, "num", len(res_genres), 4)
	})

}
