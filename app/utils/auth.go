package utils

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"example.com/app/conf"
	"github.com/golang-jwt/jwt/v4"
)

type MyCustomClaims struct {
	UserId    uint   `json:"uid"`
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

func GenerateJwt(userId uint, tokenType string, duration time.Duration) (string, error) {
	//fmt.Println("@@@@@GenerateJwt")
	// Create the claims
	claims := MyCustomClaims{
		userId,
		tokenType,
		jwt.RegisteredClaims{
			// A usual scenario is to set the expiration time relative to the current time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    conf.BACKEND_URL,
			//Subject:   "somebody",
			//ID:        "1",
			Audience: []string{conf.FRONTEND_URL},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(conf.HMAC_SECRET_KEY))
	if err != nil {
		return "", err
	}
	return ss, nil
}
func ParseJwt(tokenString string) (*MyCustomClaims, error) {
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.HMAC_SECRET_KEY), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error at ParseWithClaims: %v", err)
	}

	claims, ok := token.Claims.(*MyCustomClaims)
	if ok && token.Valid {
		//add more verifications
		//fmt.Printf("%v %v", claims.UserId, claims.RegisteredClaims.Issuer)
		if claims.Issuer != conf.BACKEND_URL {
			return nil, errors.New("invalid token issuer")
		}
		if claims.Audience[0] != conf.FRONTEND_URL {
			return nil, errors.New("invalid token audience")
		}
		fmt.Println("OK")
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
func Authenticate(authHeader, tokenType string) (*MyCustomClaims, error) {
	if authHeader == "" {
		return nil, errors.New("authorization not set")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid auth header")
	}
	token := headerParts[1]
	claim, err := ParseJwt(token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//verify claim
	if claim.TokenType != tokenType {
		return nil, errors.New("TokenType mismatch")
	}
	return claim, nil
}
