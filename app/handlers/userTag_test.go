package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/models"
	"example.com/app/utils"
)

func Test_SaveTags(t *testing.T) {

}
func Test_GetTags(t *testing.T) {

	t.Run("can get tags", func(t *testing.T) {
		h.DB = TestDB.Begin()
		defer h.DB.Rollback()

		data, err := models.InsertTestData(h.DB)
		if err != nil {
			t.Error(err)
		}
		req := httptest.NewRequest(http.MethodGet, ts.URL+"/tags", nil)
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
		var res_tags []models.UserTag
		if err := utils.BodyToStruct(res.Body, &res_tags); err != nil {
			t.Error(err)
		}
		//utils.PrintStruct(res_tags)
		got_num := len(res_tags)
		want_num := 3
		//utils.PrintStruct(res_genres)
		if got_num != want_num {
			t.Errorf("status_code: got %d, want %d", got_num, want_num)
		}
	})

}
