package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
	"github.com/google/uuid"
)

func Test_RefreshHandler(t *testing.T) {
	users, _ := models.PrepareTestUsersOnly(h.DB)
	refreshToken, _ := users[0].GenerateToken("refresh")
	sessionString := uuid.NewString()
	session := models.Session{
		UserId:        users[0].ID,
		SessionString: sessionString,
		RefreshToken:  "Bearer " + refreshToken,
	}
	session.Update(h.DB)

	//テスト後のclean up
	defer models.ClearTestDB(h.DB)
	//テストデータを定義する
	tests := []struct {
		name          string
		sessionString string
		status        int
		errorCode     int
	}{
		{"fail with no cookie", "", http.StatusBadRequest, customError.Others.Code},
		{"fail with wrong session_string", "falsesessingstring", http.StatusBadRequest, customError.Others.Code},
		{"success with valid session_string", session.SessionString, http.StatusOK, -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			//request
			cookie := utils.GetSessionCookie(test.sessionString, 10*time.Minute)
			r := httptest.NewRequest(http.MethodPost, "/test", nil)
			if test.sessionString != "" {
				r.AddCookie(cookie)
			}
			//response
			w := httptest.NewRecorder()
			//handlerの準備
			handler := http.HandlerFunc(h.RefreshHandler)
			//実行
			handler.ServeHTTP(w, r)

			//responseがwに書き込まれるのでtestと比較

			want_status := test.status
			got_status := w.Result().StatusCode
			if got_status != want_status {
				t.Errorf("statusCode: got %d, want %d", got_status, want_status)
			}
			if want_status == http.StatusOK {
				//responseの中身を見る
				type Response = struct {
					AccessToken string `json:"access_token"`
				}
				res := Response{}
				if err := json.NewDecoder(w.Result().Body).Decode(&res); err != nil {
					t.Error(err)
				}
				if res.AccessToken == "" {
					t.Error("access token not found after successful response")
				}
				got_cookie := w.Result().Cookies()[0]
				if got_cookie.Name != cookie.Name {
					got := got_cookie.Name
					want := cookie.Name
					t.Errorf("different cookie: got %s, want: %s", got, want)
				}
				if got_cookie.Value == cookie.Value {
					t.Error("refresh token not changed after successful response")
				}
			} else {
				//返ったエラーの中身を見る
				got_e_response := customError.CustomError{} //なげぇ・・・
				if err := json.NewDecoder(w.Result().Body).Decode(&got_e_response); err != nil {
					t.Error(err)
				}
				if got_e_response.Code != test.errorCode {
					t.Errorf("error response code: got %d, want %d", got_e_response.Code, test.errorCode)
				}
			}
		})
	}

}
