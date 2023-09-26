package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/app/models"
	"example.com/app/testutil"
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
		uuid        string
		condition   models.SongSearchCond
		want_num    int
		statusCode  int
	}{
		{"can get by the same user", http.MethodPost, true, data.User.UUID, searchCond, 2, http.StatusOK},
		{"cannont get by different user", http.MethodPost, true, "10000", searchCond, 0, http.StatusBadRequest},
		{"cannont access with get method", http.MethodGet, true, data.User.UUID, searchCond, 0, http.StatusBadRequest},
		{"cannont access without authorization", http.MethodPost, false, data.User.UUID, searchCond, 0, http.StatusBadRequest},
		{"cannont access without condition", http.MethodPost, true, data.User.UUID, emptyCond, 0, http.StatusBadRequest},
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
			user := &models.User{UUID: test.uuid}
			token, _ := user.FakeGenerateToken()
			if test.requireAuth {
				testutil.AddAuthorizationHeader(req, token)
			}

			res, err := ts.Client().Do(req)
			if err != nil {
				t.Error(err)
			}
			defer res.Body.Close()

			//ステータスチェック
			testutil.Checker(t, "status_code", res.StatusCode, test.statusCode)
			if res.StatusCode == http.StatusOK {
				songs := []models.UserSong{}
				utils.BodyToStruct(res.Body, &songs)
				testutil.Checker(t, "num", len(songs), test.want_num)
			}

		})
	}
}

func Test_CreateSong(t *testing.T) {
	//テストデータを定義する
	songData := models.VariousSongData()

	tests := []struct {
		name   string
		song   models.UserSong
		status int
	}{
		{"can add complex song", songData.Complex, http.StatusOK},
		{"can add simple song", songData.Simple, http.StatusOK},
		{"can add complex song with no instruments used", songData.NoInstUsed, http.StatusOK},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()
			data, _ := models.InsertTestData(h.DB)

			//request
			js, _ := utils.ToJSON(test.song)
			req := httptest.NewRequest(http.MethodPost, ts.URL+"/song/new", strings.NewReader(js))
			req.RequestURI = ""
			token, _ := data.User.FakeGenerateToken()
			testutil.AddAuthorizationHeader(req, token)

			//response
			res, err := ts.Client().Do(req)
			if err != nil {
				t.Error(err)
			}
			//responseのチェック
			testutil.Checker(t, "status_code", res.StatusCode, test.status)

			if res.StatusCode == http.StatusOK {
				us := models.UserSong{}
				utils.BodyToStruct(res.Body, &us)

				s := models.UserSong{}
				err, isFound := s.GetByUUID(h.DB, us.UUID, true)
				if !isFound {
					t.Error("failed. added song not found.")
				}
				if err != nil {
					t.Error(err)
				}
				//utils.PrintStruct(s)
			}

		})
	}

}
func Test_UpdateSong(t *testing.T) {

}
func Test_getSong(t *testing.T)    {}
func Test_DeleteSong(t *testing.T) {}
