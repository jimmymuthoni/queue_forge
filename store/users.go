package store

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)


type User struct {
	Id uuid.UUID `db:"id"`
	Email string `db:"email"`
	HashedPassword string `db:"HashedPassword"`
}

type UserStore struct {
	db *sqlx.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: sqlx.NewDb(db, "postgres"),
	}
}

