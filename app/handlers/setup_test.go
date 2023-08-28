package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"example.com/app/customError"
	"example.com/app/models"
	"example.com/app/utils"
	"gorm.io/gorm"
)

var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEwMCwidG9rZW5fdHlwZSI6ImFjY2VzcyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMCIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjkyNTA4ODc0LCJuYmYiOjE2OTI1MDg4MTQsImlhdCI6MTY5MjUwODgxNH0.Nw5g5FYh_uiZkvOg0bhxV0nIP_Z75lYZ72xjwOArbL0"

// テスト用のサーバー
var ts *httptest.Server

var h = HandlersConf{
	DB:        nil,
	IsTesting: true,
	SendEmail: false,
}
var TestDB *gorm.DB

func TestMain(m *testing.M) {
	db, err := models.InitTestDB()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	TestDB = db
	h.DB = TestDB
	ts = httptest.NewTLSServer(h.Handlers())
	defer ts.Close()
	os.Exit(m.Run())
}

/*
handlerテストのテンプレ
*/
func template(t *testing.T) {
	//テストデータを定義する
	type SomeData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	tests := []struct {
		name      string
		somedata  SomeData
		status    int
		errorCode int
	}{
		{"test 1", SomeData{"test@test.test", "dummypassword"}, http.StatusBadRequest, customError.UserAlreadyExists.Code},
		{"test 2", SomeData{"test2@test.test", "dummypassword"}, http.StatusOK, -1},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			h.DB = TestDB.Begin()
			defer h.DB.Rollback()
			//必要なデータをテーブルに入れとく
			_, err := models.InsertTestUsersOnly(h.DB)
			if err != nil {
				t.Error(err)
			}

			//request
			js, _ := utils.ToJSON(test.somedata)
			r := httptest.NewRequest(http.MethodPost, "/test", strings.NewReader(js))
			//contextにuidを仕込む
			ctx := utils.SetUIDInContext(r.Context(), uint(100))
			//response
			w := httptest.NewRecorder()
			//handlerの準備
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
			//handler := http.HandlerFunc(h.SignUpHandler)
			//実行
			handler.ServeHTTP(w, r.WithContext(ctx))

			//responseがwに書き込まれるのでtestと比較

			want_status := test.status
			got_status := w.Result().StatusCode
			if got_status != want_status {
				t.Errorf("statusCode: got %d, want %d", got_status, want_status)
			}
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
				if got_e_response.Code != test.errorCode {
					t.Errorf("error response code: got %d, want %d", got_e_response.Code, test.errorCode)
				}
			}
		})
	}

}

/*
reqにAuthorizationヘッダ付与
*/
func addAuthorizationHeader(req *http.Request, user *models.User) error {
	authorization, err := user.GenerateToken("access")
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+authorization)
	return nil
}
