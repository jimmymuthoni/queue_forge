package store

import (
	"context"
	"os"
	"testing"

	"github.com/golang-migrate/migrate"
	"github.com/jimmymuthoni/queue_forge/config"
	"github.com/stretchr/testify/require"
)

func TestUserStore(t *testing.T){
	os.Setenv("ENV", string(config.Env_Test))
	conf, err := config.New()
	require.NoError(t, err)

	db, err := NewPostgresDb(conf)
	require.NoError(t, err)
	defer db.Close()

	//migations to the test_db
	m, err := migrate.New(
		"file//migrations",
		conf.DatabaseUrl())
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	userStore := NewUserStore(db)
	userStore.CreateUser(context.Background(), "jim@testingmail.com", "testingpassword")

}