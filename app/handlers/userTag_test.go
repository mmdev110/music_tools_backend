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

		if res.StatusCode != http.StatusOK {
			got := res.StatusCode
			want := http.StatusOK
			t.Errorf("status_code: got %d, want %d", got, want)
		}
		var res_tags []models.UserTag
		if err := utils.BodyToStruct(res.Body, &res_tags); err != nil {
			t.Error(err)
		}
		//utils.PrintStruct(res_tags)
		if len(res_tags) != 3 {
			got := len(res_tags)
			want := 3
			t.Errorf("status_code: got %d, want %d", got, want)
		}
	})

}
