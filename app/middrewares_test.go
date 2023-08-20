package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/conf"
	"example.com/app/models"
	"example.com/app/utils"
)

func Test_CORS(t *testing.T) {
	emptyHandler := http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {})
	handler := enableCORS(emptyHandler)
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
			fmt.Printf("%s header want: %s, but got: %s\n", header.name, want, got)
			if got != want {
				t.Errorf("%s header want: %s, but got: %s", header.name, want, got)
			}
		}
	}

}

func Test_requireAuth(t *testing.T) {
	var idFromContext uint
	emptyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idFromContext = utils.GetUidFromContext(r.Context())
	})
	handler := requireAuth(emptyHandler)
	u := models.User{}
	u.ID = uint(100)
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
		{"no header", "", false, 400},
		{"wrong token", "wrongjwttoken", true, 400},
		{"empty token", "", true, 400},
		{"valid token", "Bearer " + token, true, 200},
		{"valid, but no bearer", token, true, 400},
		{"expired", "Bearer " + expiredToken, true, 400},
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
		if test.wantStatus == 200 {
			//contextに入れたIDのチェック
			got := idFromContext
			if got != u.ID {
				t.Errorf("ID in context: got %d, want %d", got, u.ID)
			}
		}
	}

}
