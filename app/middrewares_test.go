package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/conf"
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
