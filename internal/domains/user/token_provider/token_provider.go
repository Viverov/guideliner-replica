package tokenprovider

import (
	"time"
)

type TokenProvider interface {
	GenerateToken(userId uint, tokenTTL time.Duration) (token string, err error)
	ValidateToken(token string) (claims *AuthClaims, err error)
}

type AuthClaims struct {
	UserID    uint   `json:"userId"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}
