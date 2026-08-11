package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/golang-jwt/jwt/v4"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


type RefreshToken struct{
	UserId uuid.UUID	 	`db:"user_id"`
	HashedToken	string		`db:"hashed_token"`
	CreatedAt	time.Time	`db:"created_at"`
	ExpiresAt	time.Time 	`db:"expires_at"`
}


type RefreshTokenStore struct {
	db *sqlx.DB
}

func NewRefreshTokenStore(db *sql.DB) *RefreshTokenStore {
	return &RefreshTokenStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

// creating refresh token and adding it to the database
func (s *RefreshTokenStore) Create(ctx context.Context, user_id uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const insert = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3)`;
	hashedToken, err := bcrypt.GenerateFromPassword([]byte(token.Raw), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("bcrypt.GenerateFromPassword: %w", err)
	}
	base64TokenHash := base64.StdEncoding.EncodeToString(hashedToken)

	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, insert, user_id, base64TokenHash); err != nil {
		return nil, fmt.Errorf("failed to create refresh token record: %w", err)
	}
	return  &refreshToken, nil
	
}