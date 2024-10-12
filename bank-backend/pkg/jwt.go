package pkg

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Define these in your configuration
const (
	AccessTokenExpiration  = 24 * time.Hour
	RefreshTokenExpiration = 7 * 24 * time.Hour
	JWTSecret              = "farhan-dwian" // In production, use an environment variable
)

type Claims struct {
	PhoneNumber string `json:"phone_number"`
	jwt.RegisteredClaims
}

func GenerateAccessTokens(phone string) (string, error) {
	// Generate Access Token
	accessTokenClaims := Claims{
		PhoneNumber: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SYN-BACKEND",
			Subject:   phone,
			ID:        uuid.New().String(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}

	return accessTokenString, nil
}

func GenerateRefreshTokens(phone string) (string, error) {
	// Generate Refresh Token
	refreshTokenClaims := Claims{
		PhoneNumber: phone,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SYN-BACKEND",
			Subject:   phone,
			ID:        uuid.New().String(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(JWTSecret))
	if err != nil {
		return "", err
	}
	return refreshTokenString, nil
}
