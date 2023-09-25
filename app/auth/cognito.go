package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"example.com/app/conf"
	"example.com/app/utils"
	"github.com/golang-jwt/jwt/v4"
)

type CognitoClaims struct {
	UserName uint   `json:"user_name"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func AuthCognito(authHeader string) (*CognitoClaims, error) {
	if authHeader == "" {
		return nil, errors.New("authorization not set")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid auth header")
	}
	token := headerParts[1]
	claim, err := ParseCognitoJwt(token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//verify claim

	return claim, nil
}
func ParseCognitoJwt(tokenString string) (*CognitoClaims, error) {
	// sample token string taken from the New example
	//tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJmb28iOiJiYXIiLCJuYmYiOjE0NDQ0Nzg0MDB9.u1riaD1rW97opCoAuRCTy4w58Br-Zk-bh7vLiRIsrpU"

	token, err := jwt.ParseWithClaims(tokenString, &CognitoClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(conf.HMAC_SECRET_KEY), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error at ParseWithClaims: %v", err)
	}

	claims, ok := token.Claims.(*CognitoClaims)
	if ok && token.Valid {
		//add more verifications
		//fmt.Printf("%v %v", claims.UserId, claims.RegisteredClaims.Issuer)
		if claims.Issuer != conf.BACKEND_URL {
			return nil, errors.New("invalid token issuer")
		}
		if claims.Audience[0] != conf.FRONTEND_URL {
			return nil, errors.New("invalid token audience")
		}
		//fmt.Println("OK")
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

type JWK struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
	Use string `json:"use"`
}
type JWKS struct {
	Keys []JWK `json:"keys"`
}

// JWT検証用の鍵のリストを取得
// https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html#amazon-cognito-user-pools-using-tokens-aws-jwt-verify
func GetJWKS(awsRegion, cognitoUserPoolId string) ([]JWK, error) {
	issuer := fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", awsRegion, cognitoUserPoolId)
	jwks_uri := issuer + "/.well-known/jwks.json"
	jwks := JWKS{}

	res, errRes := http.Get(jwks_uri)
	if errRes != nil {
		return nil, errRes
	}
	defer res.Body.Close()
	//utils.BodyToString(res.Body)
	//fmt.Println(utils.BodyToString(res.Body))
	if err := utils.BodyToStruct(res.Body, &jwks); err != nil {
		return nil, err
	}

	return jwks.Keys, nil
}

func findJWKByKid(jwks []JWK, kid string) (JWK, error) {
	for _, jwk := range jwks {
		if jwk.Kid == kid {
			return jwk, nil
		}
	}
	return JWK{}, fmt.Errorf("JWK not found for kid: %s", kid)
}
