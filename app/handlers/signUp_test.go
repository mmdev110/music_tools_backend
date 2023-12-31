package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/testutil"
	"example.com/app/utils"
)

func Test_SignUpHandler(t *testing.T) {
	type Form struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	tests := []struct {
		name      string
		data      Form
		status    int
		errorCode int
	}{
		{"cannot signup with existing confirmed email", Form{"tes@test.test", "dummypassword"}, http.StatusBadRequest, customError.UserAlreadyExists.Code},
		{"can signup with existing unconfirmed email", Form{"test2@test.test", "dummypassword"}, http.StatusOK, -1},
		{"can signup with non-existing email", Form{"test3@test.test", "dummypassword"}, http.StatusOK, -1},
		{"cannot signup with empty email", Form{"", "dummypassword"}, http.StatusBadRequest, customError.InsufficientParameters.Code},
		{"cannot signup with empty password", Form{"test3@test.test", ""}, http.StatusBadRequest, customError.InsufficientParameters.Code},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//トランザクション貼る
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()
			_, err := models.InsertTestUsersOnly(h.DB)
			if err != nil {
				t.Error(err)
			}

			js, _ := utils.ToJSON(test.data)
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(js))
			//ctx := utils.SetUIDInContext(r.Context(), user.ID)

			w := httptest.NewRecorder()

			handler := http.HandlerFunc(h.SignUpHandler)
			handler.ServeHTTP(w, r)

			want_status := test.status
			got_status := w.Result().StatusCode
			testutil.Checker(t, "status_code", got_status, want_status)

			if want_status == http.StatusOK {
				//responseの中身を見る
				got_u := models.User{}
				if err := json.NewDecoder(w.Result().Body).Decode(&got_u); err != nil {
					t.Error(err)
				}
				//utils.PrintStruct(got_u)

			} else {
				//返ったエラーの中身を見る
				got_e_response := customError.CustomError{} //なげぇ・・・
				if err := json.NewDecoder(w.Result().Body).Decode(&got_e_response); err != nil {
					t.Error(err)
				}
				testutil.Checker(t, "error_code", got_e_response.Code, test.errorCode)
			}
		})
	}

}

/*
 */
func Test_EmailConfirmationHandler(t *testing.T) {
	//cognito移行により不要になったためskip
	t.Skip()
	//テストデータを定義する
	users := models.GetTestUsers()
	token_confirmed, _ := users[0].GenerateToken("email_confirm")
	token_unconfirmed, _ := users[1].GenerateToken("email_confirm")
	token_unconfirmed_wrong_type, _ := users[1].GenerateToken("access")
	userNonExistent := models.User{ID: uint(1)}
	token_non_existent, _ := userNonExistent.GenerateToken("email_confirm")
	type Token struct {
		Token string `json:"token"`
	}
	tests := []struct {
		name      string
		token     string
		status    int
		errorCode int
	}{
		{"fail with no token", "", http.StatusBadRequest, customError.InvalidToken.Code},
		{"fail with random token", "raklasfdaklw", http.StatusBadRequest, customError.InvalidToken.Code},
		{"fail with expired token", expiredToken, http.StatusBadRequest, customError.InvalidToken.Code},
		{"fail with confirmed user", token_confirmed, http.StatusBadRequest, customError.AddressAlreadyConfirmed.Code},
		{"fail with non-existent ", token_non_existent, http.StatusBadRequest, customError.UserNotFound.Code},
		{"success with unconfirmed user", token_unconfirmed, http.StatusOK, -1},
		{"fail with unconfirmed user with wrong token type", token_unconfirmed_wrong_type, http.StatusBadRequest, customError.InvalidToken.Code},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//トランザクション貼る
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()

			_, err := models.InsertTestUsersOnly(h.DB)
			if err != nil {
				t.Error(err)
			}

			//request
			js, _ := utils.ToJSON(Token{test.token})
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(js))
			//response
			w := httptest.NewRecorder()
			//handlerの準備
			handler := http.HandlerFunc(h.EmailConfirmationHandler)
			//実行
			handler.ServeHTTP(w, r)

			//responseがwに書き込まれるのでtestと比較

			want_status := test.status
			got_status := w.Result().StatusCode
			testutil.Checker(t, "status_code", got_status, want_status)

			if want_status == http.StatusOK {
				//responseの中身を見る
				got_u := models.User{}
				if err := json.NewDecoder(w.Result().Body).Decode(&got_u); err != nil {
					t.Error(err)
				}
				if !got_u.IsConfirmed {
					t.Error("user is not confirmed after success")
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
