package testutil

import (
	"net/http"
	"testing"
)

/*
reqにAuthorizationヘッダ付与
*/
func AddAuthorizationHeader(req *http.Request, token string) error {
	req.Header.Add("Authorization", "Bearer "+token)
	return nil
}

func Checker[T comparable](t *testing.T, parameterName string, got, want T) {
	if got != want {
		t.Errorf("%s: got %v, want %v", parameterName, got, want)
	}
}
