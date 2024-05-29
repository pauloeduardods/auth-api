package jwt_verify

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	AccessToken     string   `json:"accessToken"`
	Alg             string   `json:"alg"`
	Kid             string   `json:"kid"`
	Aud             string   `json:"aud"`
	AuthTime        int64    `json:"auth_time"`
	CognitoUsername string   `json:"cognito:username"`
	UserGroups      []string `json:"cognito:groups"`
	Email           string   `json:"email"`
	EmailVerified   bool     `json:"email_verified"`
	EventID         string   `json:"event_id"`
	Exp             int64    `json:"exp"`
	Iat             int64    `json:"iat"`
	Iss             string   `json:"iss"`
	Jti             string   `json:"jti"`
	Name            string   `json:"name"`
	OriginJti       string   `json:"origin_jti"`
	Sub             string   `json:"sub"`
	TokenUse        string   `json:"token_use"`
}

func (c *Claims) GetExpirationTime() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.Exp, 0)}, nil
}

func (c *Claims) GetIssuedAt() (*jwt.NumericDate, error) {
	return &jwt.NumericDate{Time: time.Unix(c.Iat, 0)}, nil
}

func (c *Claims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c *Claims) GetIssuer() (string, error) {
	return c.Iss, nil
}

func (c *Claims) GetSubject() (string, error) {
	return c.Sub, nil
}

func (c *Claims) GetAudience() (jwt.ClaimStrings, error) {
	return jwt.ClaimStrings{c.Aud}, nil
}
