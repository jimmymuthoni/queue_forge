package store

import "database/sql"

type Store struct {
	Users *UserStore
	RefreshTokenStore *RefreshTokenStore
}

func New(db *sql.DB) *Store {
	return &Store{
		Users: NewUserStore(db),
		RefreshTokenStore: NewRefreshTokenStore(db),
	}
}