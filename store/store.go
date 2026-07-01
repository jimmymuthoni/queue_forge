package store

import "database/sql"

type Store struct {
	Users *UserStore
}

func New(db *sql.DB) *Store {
	return &Store{
		Users: NewUserStore(db),
	}
}