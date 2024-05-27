package jwt_verify

import (
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"monitoring-system/server/pkg/logger"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

type JWTVerify interface {
	CacheJWK() error
	ParseJWT(tokenString string) (*jwt.Token, error)
	JWK() *JWK
	JWKURL() string
}

type jwtVerify struct {
	jwk               *JWK
	jwkURL            string
	cognitoRegion     string
	cognitoUserPoolID string
	log               logger.Logger
}

type JWK struct {
	Keys []struct {
		Alg string `json:"alg"`
		E   string `json:"e"`
		Kid string `json:"kid"`
		Kty string `json:"kty"`
		N   string `json:"n"`
	} `json:"keys"`
}

func NewAuth(cognitoRegion, cognitoUserPoolID string, logger logger.Logger) JWTVerify {
	a := &jwtVerify{
		cognitoRegion:     cognitoRegion,
		cognitoUserPoolID: cognitoUserPoolID,
		log:               logger,
	}

	a.jwkURL = fmt.Sprintf("https://cognito-idp.%s.amazonaws.com/%s/.well-known/jwks.json", a.cognitoRegion, a.cognitoUserPoolID)

	return a
}

func (a *jwtVerify) CacheJWK() error { // Check when we need to cache the JWK
	req, err := http.NewRequest("GET", a.jwkURL, nil)
	if err != nil {
		a.log.Error("Error creating JWK request %v", err)
		return err
	}

	req.Header.Add("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		a.log.Error("Error getting JWK response %v", err)
		return err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		a.log.Error("Error reading JWK response body %v", err)
		return err
	}

	jwk := new(JWK)
	err = json.Unmarshal(body, jwk)
	if err != nil {
		a.log.Error("Error unmarshalling JWK %v", err)
		return err
	}

	a.jwk = jwk
	return nil
}

func (a *jwtVerify) ParseJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		key, err := convertKey(a.jwk.Keys[1].E, a.jwk.Keys[1].N)
		return key, err
	})
	if err != nil {
		a.log.Error("Error parsing JWT %v", err)
		return token, err
	}

	return token, nil
}

func (a *jwtVerify) JWK() *JWK {
	return a.jwk
}

func (a *jwtVerify) JWKURL() string {
	return a.jwkURL
}

func convertKey(rawE, rawN string) (*rsa.PublicKey, error) {
	decodedE, err := base64.RawURLEncoding.DecodeString(rawE)
	if err != nil {
		return nil, err
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
	decodedN, err := base64.RawURLEncoding.DecodeString(rawN)
	if err != nil {
		return nil, err
	}
	pubKey.N.SetBytes(decodedN)
	return pubKey, nil
}
