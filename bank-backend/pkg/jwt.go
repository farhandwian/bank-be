package pkg

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

// Define these in your configuration
const (
	AccessTokenExpiration  = 5 * time.Minute
	RefreshTokenExpiration = 7 * 24 * time.Hour
	JWTSecret              = "harisatul-aulia" // In production, use an environment variable
)

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

func GenerateAccessTokens(username string, role string) (string, error) {
	// Generate Access Token
	accessTokenClaims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SYN-BACKEND",
			Subject:   username,
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

func GenerateRefreshTokens(username string) (string, error) {
	// Generate Refresh Token
	refreshTokenClaims := Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenExpiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "SYN-BACKEND",
			Subject:   username,
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
