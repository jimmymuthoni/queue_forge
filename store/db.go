package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/jimmymuthoni/queue_forge/config"
	_ "github.com/lib/pq"
)


//conecting to th database
func NewPostgresDb(conf *config.Config) (*sql.DB, error){
	dsn := conf.DatabaseUrl()
	db,err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open database connection: %w", err)
	}
	ctx,cancel := context.WithTimeout(context.Background(),5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping the database")
	}
	return db, nil
}