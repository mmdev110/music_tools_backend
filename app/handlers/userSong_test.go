package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/app/models"
	"example.com/app/utils"
)

func Test_SearchSongsHandler(t *testing.T) {
	data := models.GetTestData()

	searchCond := models.SongSearchCond{
		UserIds: []uint{data.User.ID},
	}
	emptyCond := models.SongSearchCond{}
	tests := []struct {
		name        string
		method      string
		requireAuth bool
		uid         uint
		condition   models.SongSearchCond
		want_num    int
		statusCode  int
	}{
		{"can get by the same user", http.MethodPost, true, data.User.ID, searchCond, 2, http.StatusOK},
		{"cannont get by different user", http.MethodPost, true, uint(10000), searchCond, 0, http.StatusBadRequest},
		{"cannont access with get method", http.MethodGet, true, data.User.ID, searchCond, 0, http.StatusBadRequest},
		{"cannont access without authorization", http.MethodPost, false, data.User.ID, searchCond, 0, http.StatusBadRequest},
		{"cannont access without condition", http.MethodPost, true, data.User.ID, emptyCond, 0, http.StatusBadRequest},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()

			models.InsertTestUsersOnly(h.DB)
			models.InsertTestData(h.DB)

			js, err := utils.ToJSON(test.condition)
			req := httptest.NewRequest(test.method, ts.URL+"/list", strings.NewReader(js))
			req.RequestURI = ""
			if test.requireAuth {
				addAuthorizationHeader(req, &models.User{ID: test.uid})
			}

			res, err := ts.Client().Do(req)
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()

			//ステータスチェック
			checker(t, "status_code", res.StatusCode, test.statusCode)
			if res.StatusCode == http.StatusOK {
				songs := []models.UserSong{}
				utils.BodyToStruct(res.Body, &songs)
				checker(t, "num", len(songs), test.want_num)
			}

		})
	}
}

func Test_CreateSong(t *testing.T) {}
func Test_UpdateSong(t *testing.T) {

}
func Test_getSong(t *testing.T)    {}
func Test_DeleteSong(t *testing.T) {}
