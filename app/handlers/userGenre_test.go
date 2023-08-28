package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/models"
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
		addAuthorizationHeader(req, data.User)

		res, err := ts.Client().Do(req)
		if err != nil {
			t.Error(err)
		}
		defer res.Body.Close()

		got_code := res.StatusCode
		want_code := http.StatusOK
		if got_code != want_code {
			t.Errorf("status_code: got %d, want %d", got_code, want_code)
		}
		var res_genres []models.UserGenre
		if err := utils.BodyToStruct(res.Body, &res_genres); err != nil {
			t.Error(err)
		}

		got_num := len(res_genres)
		want_num := 4
		//utils.PrintStruct(res_genres)
		if got_num != want_num {
			t.Errorf("status_code: got %d, want %d", got_num, want_num)
		}
	})

}
