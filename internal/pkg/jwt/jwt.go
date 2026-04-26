package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token expired")
)

type Claims struct {
	UserID  string `json:"user_id"`
	TokenID string `json:"token_id"`
	jwt.RegisteredClaims
}

type JWT struct {
	secret            string
	accessExpiry      time.Duration
	refreshExpiry    time.Duration
}

func NewJWT(secret string, accessExpiry, refreshExpiry time.Duration) *JWT {
	return &JWT{
		secret:       secret,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

func (j *JWT) GenerateAccessToken(userID string) (string, error) {
	tokenID, err := generateTokenID()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID:  userID,
		TokenID: tokenID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.secret))
}

func (j *JWT) ValidateAccessToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, ErrExpiredToken
	}

	return claims, nil
}

func (j *JWT) GenerateRefreshToken() (string, error) {
	return generateTokenID()
}

func (j *JWT) GetRefreshTokenExpiry() time.Duration {
	return j.refreshExpiry
}

func generateTokenID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.RawStdEncoding.EncodeToString(b), nil
}