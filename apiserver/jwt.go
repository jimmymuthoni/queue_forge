package apiserver

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jimmymuthoni/queue_forge/config"
)

type JwtManager struct {
	config *config.Config
}

//storing the tokens
type TokenPair struct {
	AccessToken  *jwt.Token
	RefreshToken *jwt.Token
}


func NewJwtManager(config *config.Config) *JwtManager {
	return &JwtManager{config}
}

type CustomClaims struct {
	TokenType string `json:"token_type"`
	jwt.RegisteredClaims
}

//this function genartes new tokens and store them in token pair struct
func (j *JwtManager) GenerateTokenPair(userId uuid.UUID) (*TokenPair, error){
	now := time.Now()
	issuer := "http://" + j.config.ApiServerHost + ":" + j.config.ApiServerPort

	//geting key for access token
	jwtAccessToken := jwt.NewWithClaims(jwt.SigningMethodES256, CustomClaims {
		TokenType: "access",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userId.String(),
			Issuer: issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour * 24 * 30)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	})
	key := []byte(j.config.JwtSecret)
	var err error
	jwtAccessToken.Raw, err = jwtAccessToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	//key for refresh token
	jwtRefreshToken := jwt.NewWithClaims(jwt.SigningMethodES256, CustomClaims {
		TokenType: "refresh",
		RegisteredClaims: jwt.RegisteredClaims{
			Subject: userId.String(),
			Issuer: issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Minute * 15)),
			IssuedAt: jwt.NewNumericDate(now),
		},
	})

	jwtRefreshToken.Raw, err = jwtRefreshToken.SignedString(key)
	if err != nil {
		return nil, fmt.Errorf("failed to sign refresh token: %w", err)
	}
	return &TokenPair{
		AccessToken: jwtAccessToken,
		RefreshToken: jwtRefreshToken,
	}, nil
}