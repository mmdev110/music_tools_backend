package auth

import (
	"testing"

	"example.com/app/utils"
)

func TestAuthConfigure(t *testing.T) {
	t.Run("test auth.Configure()", func(t *testing.T) {
		utils.PrintStruct(auth)
		for _, key := range auth.JWKS {
			utils.PrintStruct(key)
		}
	})

}
func TestAuthCognito(t *testing.T) {
	t.Skip()
	validToken := "Bearer JWT_TOKEN"
	t.Run("AuthCognito", func(t *testing.T) {
		claims, err := auth.AuthCognito(validToken)
		if err != nil {
			//t.Error(err)
		}
		utils.PrintStruct(claims)
	})

}
