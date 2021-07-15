// +build unit

package token_provider

import (
	"github.com/bxcodec/faker/v3"
	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// Should create token with userId without errors
func TestGenerateToken(t *testing.T) {
	jwtService, _, userId := setupTestData()

	token, err := jwtService.GenerateToken(userId, time.Hour)

	assert.Nil(t, err)
	assert.NotEqual(t, "", token)
}

// Should validate token, created by GenerateToken
func TestValidateToken(t *testing.T) {
	jwtService, _, userId := setupTestData()

	token, _ := jwtService.GenerateToken(userId, time.Hour)
	claims, err := jwtService.ValidateToken(token)

	assert.Nil(t, err)
	assert.Equal(t, claims.UserId, userId)
}

func TestValidateTokenExpired(t *testing.T) {
	jwtService, secretKey, _ := setupTestData()

	token := generateExpiredToken(secretKey)
	claims, err := jwtService.ValidateToken(token)

	assert.Nil(t, claims)
	assert.NotNil(t, err)
	assert.EqualError(t, err, (&ExpiredTokenError{}).Error())
}

func TestValidateInvalidToken(t *testing.T) {
	jwtService, _, _ := setupTestData()

	token := faker.Word()

	claims, err := jwtService.ValidateToken(token)

	assert.Nil(t, claims)
	assert.NotNil(t, err)
	assert.EqualError(t, err, (&NotTokenError{token: token}).Error())
}

func setupTestData() (service TokenProvider, secretKey string, userId uint) {
	secretKey = faker.Word()
	numbers, _ := faker.RandomInt(1, 1000)
	userId = uint(numbers[0])
	return NewTokenServiceJWT(secretKey), secretKey, userId
}

func generateExpiredToken(secretKey string) string {
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(-1)).Unix(),
		Issuer:    "expired",
		IssuedAt:  time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secretKey))
	return tokenString
}
