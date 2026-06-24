package store

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id                   uuid.UUID `db:"id"`
	Email                string    `db:"email"`
	HashedPasswordBase64 string    `db:"hashed_password"`
	CreatedAt 			 time.Time `db:"created_at"`
}

func (u *User) ComparePassword(password string) error {
	hashedPassword, err := base64.StdEncoding.DecodeString(u.HashedPasswordBase64)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
	if err != nil{
		return fmt.Errorf("invalid password")
	}
	return nil
}

type UserStore struct {
	db *sqlx.DB
}

//conveting sql.Db to sqlx
func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

func (s *UserStore) CreateUser(ctx context.Context, email, password string) (*User, error){
	const dml = `INSERT INTO users (email,hashed_password) VALUES ($1, $2) RETURNING *`
	var user User

	//hashing password provided by the user before inserting it to db
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Failed to hash the password: %w", err)
	}
	hashedPasswordBase64 := base64.StdEncoding.EncodeToString(bytes)

	if err := s.db.GetContext(ctx, &user, dml, email, hashedPasswordBase64); err != nil {
		return nil, fmt.Errorf("Failed to insert the user")
	}
	return &user, nil
}

//finding user by email
func (s *UserStore) FindUserByEmail(ctx context.Context, userEmail string) (*User, error){
	const query = "SELECT * FROM users WHERE email=$1;"
	var user User

	if err := s.db.GetContext(ctx, &user, query, userEmail); err != nil {
		return nil, fmt.Errorf("Failed to fetch the user: %w", err)
	}
	return &user, nil
}

//finding user by id
func (s *UserStore) FindUserById(ctx context.Context, userId uuid.UUID) (*User, error){
	const query = "SELECT * FROM users WHERE id=$1;"
	var user User

	if err := s.db.GetContext(ctx, &user, query, userId); err != nil {
		return nil, fmt.Errorf("Failed to fetch the user by Id %s: %w", userId, err)
	}
	return &user, nil
}

