package utils

import (
	"fmt"
	"net/http"
	"testing"
	"time"
)

func TestPresignedURL(t *testing.T) {
	t.Run("test presignedURL", func(t *testing.T) {
		path := "1/"
		duration := time.Duration(10 * time.Minute)
		url, err := GenerateSignedUrl(path, http.MethodGet, duration)
		if err != nil {
			t.Error(err)
		}
		fmt.Println(url)
	})

}
