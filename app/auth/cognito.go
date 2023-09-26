package auth

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	"example.com/app/utils"
	"github.com/golang-jwt/jwt/v4"
)

type AuthFunc func(authHeader string) (*CognitoClaims, error)

// jwtの中身
type CognitoClaims struct {
	UUID     string `json:"cognito:username"`
	Email    string `json:"email"`
	TokenUse string `json:"token_use"`
	jwt.RegisteredClaims
}

// JWT検証用の公開鍵
// https://docs.aws.amazon.com/ja_jp/cognito/latest/developerguide/amazon-cognito-user-pools-using-tokens-verifying-a-jwt.html#amazon-cognito-user-pools-using-tokens-aws-jwt-verify
type JWK struct {
	Kid string `json:"kid"`
	Alg string `json:"alg"`
	Kty string `json:"kty"`
	E   string `json:"e"`
	N   string `json:"n"`
	Use string `json:"use"`
}
type Auth struct {
	JWKS        []JWK `json:"keys"`
	AwsRegion   string
	UserPoolID  string
	AppClientID string
	Issuer      string
	configured  bool
}

func (auth *Auth) AuthCognito(authHeader string) (*CognitoClaims, error) {
	if err := auth.requireConfiguration(); err != nil {
		return nil, err
	}
	if authHeader == "" {
		return nil, errors.New("authorization not set")
	}
	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return nil, errors.New("invalid auth header")
	}
	token := headerParts[1]
	claim, err := auth.parseCognitoJwt(token)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	//verify claim

	return claim, nil
}
func (auth *Auth) parseCognitoJwt(tokenString string) (*CognitoClaims, error) {
	if err := auth.requireConfiguration(); err != nil {
		return nil, err
	}
	token, err := jwt.ParseWithClaims(tokenString, &CognitoClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		//check alg
		//utils.PrintStruct(token)
		if token.Method.Alg() != jwt.SigningMethodRS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		//find kid from public keys
		kid, _ := token.Header["kid"].(string)
		jwk, err := auth.findJWKByKid(kid)
		if err != nil {
			//kid見つからなかった場合はcognito側でキーがローテーションされ、
			//バックエンドで保持しているauth.JWKSが古くなった可能性があるので、auth.JWKSを取得し直す必要がある
			//バックエンド再起動または定期的なauth.JWKS再取得については後回し
			return nil, err
		}
		return jwk.convertKey(), nil
	})
	if err != nil {
		//署名ミスマッチ、期限切れはここでキャッチされる
		return nil, fmt.Errorf("error at ParseWithClaims: %v", err)
	}
	claims, ok := token.Claims.(*CognitoClaims)
	if ok && token.Valid {
		//jwtが改竄されてない&期限内であることがわかったので
		//中身(claims)をチェック
		//「クレームを検証する」参照

		//Issuer(iss)
		if claims.Issuer != auth.Issuer {
			return nil, errors.New("invalid token issuer")
		}
		//Audience(aud)
		if claims.Audience[0] != auth.AppClientID {
			return nil, errors.New("invalid token audience")
		}
		//token_use
		if claims.TokenUse != "id" {
			return nil, errors.New("invalid token_use")
		}
		//fmt.Println("OK")
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

// JWT検証用の鍵のリストを取得
func (auth *Auth) getJWKS() ([]JWK, error) {
	if auth.Issuer == "" {
		return nil, errors.New("Issuer not provided")
	}
	jwks_uri := auth.Issuer + "/.well-known/jwks.json"

	res, errRes := http.Get(jwks_uri)
	if errRes != nil {
		return nil, errRes
	}
	defer res.Body.Close()
	//utils.BodyToString(res.Body)
	//fmt.Println(utils.BodyToString(res.Body))
	a := Auth{}
	if err := utils.BodyToStruct(res.Body, &a); err != nil {
		return nil, err
	}

	return a.JWKS, nil
}

// 鍵リストからkidにマッチしたものを探す
func (auth *Auth) findJWKByKid(kid string) (*JWK, error) {
	if err := auth.requireConfiguration(); err != nil {
		return nil, err
	}
	for _, jwk := range auth.JWKS {
		if jwk.Kid == kid {
			return &jwk, nil
		}
	}
	return nil, fmt.Errorf("JWK not found for kid: %s", kid)
}

func (auth *Auth) Configure(region, userPoolID, appClientID string) error {
	auth.AwsRegion = region
	auth.UserPoolID = userPoolID
	auth.AppClientID = appClientID
	auth.Issuer = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s", auth.AwsRegion, auth.UserPoolID)
	jwks, err := auth.getJWKS()
	if err != nil {
		return err
	}
	auth.JWKS = jwks
	auth.configured = true
	return nil
}

func (auth *Auth) requireConfiguration() error {
	if !auth.configured {
		return errors.New("must be configured before running this function")
	}
	return nil
}

// https://gist.github.com/miguelmota/06f563756448b0d4ce2ba508b3cbe6e2
func (jwk *JWK) convertKey() *rsa.PublicKey {
	decodedE, err := base64.RawURLEncoding.DecodeString(jwk.E)
	if err != nil {
		panic(err)
	}
	if len(decodedE) < 4 {
		ndata := make([]byte, 4)
		copy(ndata[4-len(decodedE):], decodedE)
		decodedE = ndata
	}
	pubKey := &rsa.PublicKey{
		N: &big.Int{},
		E: int(binary.BigEndian.Uint32(decodedE[:])),
	}
	decodedN, err := base64.RawURLEncoding.DecodeString(jwk.N)
	if err != nil {
		panic(err)
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey
}
