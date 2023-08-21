package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/utils"
)

func Test_Chk(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/_chk", nil)
	handler := http.HandlerFunc(ChkHandler)

	want_status := http.StatusOK

	mp := map[string]string{"Status": "OK"}

	want_response, err := utils.ToJSON(mp)
	if err != nil {
		t.Error(err)
	}
	handler.ServeHTTP(w, r)

	got_status := w.Result().StatusCode
	got_response := utils.BodyToString(w.Result().Body)

	if got_status != want_status {
		t.Errorf("status: got %d, want %d", got_status, want_status)
	}
	if got_response != want_response {
		t.Errorf("response: got %s, want %s", got_response, want_response)
	}
}
