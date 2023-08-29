package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/customError"
	"example.com/app/testutil"
)

func Test_ResponseJSON(t *testing.T) {
	type Payload1 struct {
		Message string `json:"message"`
	}
	type DummySong struct {
		ID     uint   `json:"id"`
		UUID   string `json:"uuid"`
		UserId uint   `json:"user_id"`
		Title  string `json:"title"`
	}
	type DummyUser struct {
		ID          uint        `json:"id"`
		Email       string      `json:"email"`
		IsConfirmed bool        `json:"is_confirmed"`
		Password    string      `json:"-"`
		Songs       []DummySong `json:"songs"`
	}

	tests := []struct {
		name    string
		payload interface{}
		status  int
	}{
		{name: "messagejson", status: http.StatusOK, payload: Payload1{Message: "hi!"}},
		{name: "bad status", status: http.StatusBadRequest, payload: Payload1{Message: "hi!"}},
		{name: "complex", status: http.StatusOK, payload: DummyUser{
			ID:          uint(100),
			Email:       "aaa@aaaa.aaaa",
			IsConfirmed: false,
			Password:    "dummypassword",
			Songs: []DummySong{
				{ID: uint(1), UUID: "1jieeieweiwj", UserId: uint(100), Title: "dummytitle"},
			},
		}},
	}

	for _, test := range tests {
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ResponseJSON(w, test.payload, test.status)
		})
		r := httptest.NewRequest(http.MethodGet, "http://test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		if w.Result().StatusCode != test.status {
			t.Errorf("StatusCode: got %d, want %d", w.Result().StatusCode, test.status)
		}
		got_body := BodyToString(w.Result().Body)
		b, err := json.MarshalIndent(test.payload, "", "\t")
		want_body := string(b)
		if err != nil {
			t.Error(err)
		}
		testutil.Checker(t, "body_string", got_body, want_body)

	}
}

func Test_ErrorJSON(t *testing.T) {
	tests := []struct {
		name          string
		err           customError.CustomError
		additionalErr error
	}{
		{name: "usernotfound", err: customError.UserNotFound, additionalErr: nil},
		{name: "usernotfound", err: customError.UserNotFound, additionalErr: errors.New("lalalala")},
	}

	for _, test := range tests {
		fmt.Println(test.name)
		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ErrorJSON(w, test.err, test.additionalErr)
		})
		r := httptest.NewRequest(http.MethodGet, "http://test", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, r)

		got_status := w.Result().StatusCode
		want_status := http.StatusBadRequest
		testutil.Checker(t, "status_code", got_status, want_status)

		got_body := BodyToString(w.Result().Body)
		payload := test.err
		if test.additionalErr != nil {
			payload.Message = test.additionalErr.Error()
		}
		b, err := json.MarshalIndent(payload, "", "\t")
		want_body := string(b)
		if err != nil {
			t.Error(err)
		}
		testutil.Checker(t, "body_string", got_body, want_body)

	}
}
