package store

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

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


func (s *RefreshTokenStore) Create(ctx context.Context, token *jwt.Token) (*RefreshToken, error) {
	const insert = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3)`;
}