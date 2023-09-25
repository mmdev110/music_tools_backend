package auth

import (
	"testing"

	"example.com/app/utils"
)

func TestGetJWKS(t *testing.T) {
	jwks, err := GetJWKS("ap-northeast-1", "ap-northeast-1_05yjCtNG1")
	if err != nil {
		t.Error(err)
	}
	for _, key := range jwks {
		utils.PrintStruct(key)
	}
}
