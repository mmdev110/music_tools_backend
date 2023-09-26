package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/auth"
	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/testutil"
	"example.com/app/utils"
)

var expiredToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1aWQiOjEwMCwidG9rZW5fdHlwZSI6ImFjY2VzcyIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6NTAwMCIsImF1ZCI6WyJodHRwOi8vbG9jYWxob3N0OjMwMDAiXSwiZXhwIjoxNjkyNTA4ODc0LCJuYmYiOjE2OTI1MDg4MTQsImlhdCI6MTY5MjUwODgxNH0.Nw5g5FYh_uiZkvOg0bhxV0nIP_Z75lYZ72xjwOArbL0"

func Test_CORS(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {})
	handler := EnableCORS(emptyHandler)
	type Header = struct {
		name  string
		value string
	}
	tests := []struct {
		name    string
		method  string
		headers []Header
	}{
		{"get", http.MethodGet, []Header{
			{"Access-Control-Allow-Origin", conf.FRONTEND_URL},
			{"Access-Control-Allow-Credentials", "true"},
		}},
		{"preflight", http.MethodOptions, []Header{
			{"Access-Control-Allow-Origin", conf.FRONTEND_URL},
			{"Access-Control-Allow-Credentials", "true"},
			{"Access-Control-Allow-Headers", "Accept, Content-Type, X-CSRF-Token, Authorization"},
			{"Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS"},
		}},
	}
	for _, test := range tests {
		r := httptest.NewRequest(test.method, "http://test/", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		for _, header := range test.headers {
			got := w.Header().Get(header.name)
			want := header.value
			testutil.Checker(t, "header_value", got, want)

		}
	}

}

func Test_requireAuth(t *testing.T) {
	//cognitoと結合してテストが困難なので一旦テストを止める
	t.Skip()
	var uuid_got string
	var email_got string
	_ = uuid_got
	_ = email_got

	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		uuid_got, email_got = utils.GetParamsFromContext(r.Context())
	})

	auth := &auth.Auth{}
	handler := RequireAuth(emptyHandler, auth.AuthCognito)
	u := models.User{
		ID:   uint(100),
		UUID: "100",
	}
	token, err := u.GenerateToken("access")

	if err != nil {
		t.Error(err)
	}

	tests := []struct {
		name       string
		token      string
		setHeader  bool
		wantStatus int
	}{

		{"no header", "", false, http.StatusBadRequest},
		{"wrong token", "wrongjwttoken", true, http.StatusBadRequest},
		{"empty token", "", true, http.StatusBadRequest},
		{"valid token", "Bearer " + token, true, http.StatusOK},
		{"valid, but no bearer", token, true, http.StatusBadRequest},
		{"expired", "Bearer " + expiredToken, true, http.StatusBadRequest},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "http://test", nil)
		if test.setHeader {
			r.Header.Add("Authorization", test.token)
		}
		handler.ServeHTTP(w, r)
		gotStatus := w.Result().StatusCode
		if gotStatus != test.wantStatus {
			t.Errorf("status: got %d, want %d", gotStatus, test.wantStatus)
		}
		if test.wantStatus == http.StatusOK {
			//contextに入れたIDのチェック
			got := uuid_got
			if got != u.UUID {
				t.Errorf("ID in context: got %s, want %s", got, u.UUID)
			}
		}
	}
}
