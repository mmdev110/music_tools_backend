package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/app/utils"
)

func Test_Chk(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/_chk", nil)
	handler := http.HandlerFunc(h.ChkHandler)

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

func Test_ChkServer(t *testing.T) {
	fmt.Println(ts.URL)
	r := httptest.NewRequest(http.MethodGet, ts.URL+"/_chk", nil)
	fmt.Println(r.RequestURI)
	fmt.Println(r.URL)
	r.RequestURI = ""
	w, err := ts.Client().Do(r)
	//w, err := ts.Client().Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	defer w.Body.Close()

	if w.StatusCode != http.StatusOK {
		t.Errorf("aaaaaa%d", w.StatusCode)
	}
	fmt.Println(utils.BodyToString(w.Body))
}
