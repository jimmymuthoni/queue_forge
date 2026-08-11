package store

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/base64"
	"fmt"
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

//converting  bytes of tokens into strings
func (s *RefreshTokenStore) getBase64HashFromToken(token *jwt.Token) (string, error) {
	h := sha256.New()
	h.Write([]byte(token.Raw))
	hashedBytes := h.Sum(nil)
	base64TokenHash := base64.StdEncoding.EncodeToString(hashedBytes)
	return base64TokenHash, nil
}

// creating refresh token and adding it to the database
func (s *RefreshTokenStore) Create(ctx context.Context, user_id uuid.UUID, token *jwt.Token) (*RefreshToken, error) {
	const insert = `INSERT INTO refresh_tokens (user_id, hashed_token, expires_at) VALUES ($1, $2, $3) RETURNING *`;
	base64TokenHash, err := s.getBase64HashFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get base 64 encode token hash: %w", err)
	}
	
	expiresAt, err := token.Claims.GetExpirationTime()
	if err != nil {
		return nil, fmt.Errorf("failed to extract expiration time: %w", err)
	}
	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, insert, user_id, base64TokenHash, expiresAt.Time); err != nil {
		return nil, fmt.Errorf("failed to create refresh token record: %w", err)
	}
	return  &refreshToken, nil	
}

func (s *RefreshTokenStore) ByPrimaryKey(ctx context.Context, userId uuid.UUID, token *jwt.Token) (*RefreshToken, error){
	const query = "SELECT * FROM refresh_tokens WHERE user_id = $1 AND hashed_token = $2";
	base64TokenHash, err := s.getBase64HashFromToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to get base 64 encode token hash: %w", err)
	}
	var refreshToken RefreshToken
	if err := s.db.GetContext(ctx, &refreshToken, query, userId, base64TokenHash); err != nil {
		return nil, fmt.Errorf("failed to fetch hashed_token %s record for user %s: %w", base64TokenHash)
	}
	return &refreshToken, nil

}

func (s *RefreshTokenStore) DeleteUserTokens(ctx context.Context, userId uuid.UUID) (sql.Result, error) {
	const deleteStatement = `DELETE FROM refresh_token WHERE user_id = $1 AND hashed_token= $2`
	result, err := s.db.ExecContext(ctx, deleteStatement, userId)
	if err != nil {
		return result, fmt.Errorf("failed to delete refresh_token recordL %w", err)
	}
	return result, nil
}
