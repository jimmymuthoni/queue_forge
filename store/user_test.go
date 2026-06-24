package store

import (
	"context"
	"os"
	"testing"

	migrate "github.com/golang-migrate/migrate/v4"
    _ "github.com/golang-migrate/migrate/v4/database/postgres"
    _ "github.com/golang-migrate/migrate/v4/source/file"
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
		"file://../migrations",
		conf.DatabaseUrl())
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}

	userStore := NewUserStore(db)
	user, err := userStore.CreateUser(context.Background(), "jim@testingmail.com", "testingpassword")
	require.NoError(t, err)

	require.Equal(t, "jim@testingmail.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))


}