package token_provider

import (
	"github.com/dgrijalva/jwt-go"
	"time"
)

type tokenProviderJWT struct {
	secretKey string
	issure    string
}

type authJwtClaim struct {
	UserId uint
	jwt.StandardClaims
}

func NewTokenServiceJWT(secretKey string) TokenProvider {
	return &tokenProviderJWT{
		secretKey: secretKey,
		issure:    "guideliner",
	}
}

func (s *tokenProviderJWT) GenerateToken(userId uint, tokenTTL time.Duration) (string, error) {
	claims := &authJwtClaim{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			Issuer:    s.issure,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//encoded string
	t, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", &UnexpectedGenerateError{}
	}
	return t, nil
}

func (s *tokenProviderJWT) ValidateToken(encodedToken string) (*AuthClaims, error) {
	token, err := jwt.ParseWithClaims(encodedToken, &authJwtClaim{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				return nil, &NotTokenError{token: encodedToken}
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				return nil, &ExpiredTokenError{}
			}
		}

		return nil, &UnexpectedTokenError{token: encodedToken}
	}

	jwtClaims := token.Claims.(*authJwtClaim)
	return &AuthClaims{
		UserId:    jwtClaims.UserId,
		Audience:  jwtClaims.StandardClaims.Audience,
		ExpiresAt: jwtClaims.StandardClaims.ExpiresAt,
		Id:        jwtClaims.StandardClaims.Id,
		IssuedAt:  jwtClaims.StandardClaims.IssuedAt,
		Issuer:    jwtClaims.StandardClaims.Issuer,
		NotBefore: jwtClaims.StandardClaims.NotBefore,
		Subject:   jwtClaims.StandardClaims.Subject,
	}, nil
}
