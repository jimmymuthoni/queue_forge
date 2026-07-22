package apiserver

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jimmymuthoni/queue_forge/config"
)

var signingMethod = jwt.SigningMethodHS256

type JwtManager struct {
	config *config.Config
}

type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}

func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{
		config: config,
	}
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

// Parse verifies the JWT and returns it.
func (j *JwtManager) Parse(tokenString string) (*jwt.Token, error) {
	parser := jwt.NewParser()

	token, err := parser.ParseWithClaims(
		tokenString,
		&CustomClaims{},
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != signingMethod {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}

			return []byte(j.config.JwtSecret), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	return token, nil
}

// IsAccessToken returns true if the token is an access token.
func (j *JwtManager) IsAccessToken(token *jwt.Token) bool {
	claims, ok := token.Claims.(CustomClaims)
	if !ok {
		if c, ok := token.Claims.(*CustomClaims); ok {
			return c.TokenType == "access"
		}
		return false
	}

	return claims.TokenType == "access"
}

// GenerateTokenPair creates a new access and refresh token.
func (j *JwtManager) GenerateTokenPair(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()

	issuer := "http://" + j.config.ApiServerHost + ":" + j.config.ApiServerPort
	key := []byte(j.config.JwtSecret)

	// Access Token (15 minutes)
	accessToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(15 * time.Minute)),
		},
	})

	var err error
	accessToken.Raw, err = accessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	// Refresh Token (30 days)
	refreshToken := jwt.NewWithClaims(signingMethod, CustomClaims{
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			Issuer:    issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(30 * 24 * time.Hour)),
		},
	})

	refreshToken.Raw, err = refreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}