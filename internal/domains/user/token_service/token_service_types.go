package token_service

import (
	"fmt"
	"time"
)

type TokenServicer interface {
	GenerateToken(userId uint, tokenTTL time.Duration) (token string, err error)
	ValidateToken(token string) (claims *AuthClaims, err error)
}

type AuthClaims struct {
	UserId    uint   `json:"userId"`
	Audience  string `json:"aud,omitempty"`
	ExpiresAt int64  `json:"exp,omitempty"`
	Id        string `json:"jti,omitempty"`
	IssuedAt  int64  `json:"iat,omitempty"`
	Issuer    string `json:"iss,omitempty"`
	NotBefore int64  `json:"nbf,omitempty"`
	Subject   string `json:"sub,omitempty"`
}

type NotTokenError struct {
	token string
}

func (e *NotTokenError) Error() string {
	return fmt.Sprintf("Provided string is not a token: %s", e.token)
}

type ExpiredTokenError struct{}

func (e *ExpiredTokenError) Error() string {
	return "Provided token is expired"
}

type UnexpectedTokenError struct {
	token string
}

func (e *UnexpectedTokenError) Error() string {
	return fmt.Sprintf("Can't handle provided token: %s", e.token)
}
