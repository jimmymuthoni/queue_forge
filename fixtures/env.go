package fixtures

import (
	"database/sql"
	"os"
	"testing"
	migrate "github.com/golang-migrate/migrate/v4"

	"github.com/jimmymuthoni/queue_forge/config"
	"github.com/jimmymuthoni/queue_forge/store"
	"github.com/stretchr/testify/require"
)

type TestEnv struct {
	Config  *config.Config
	Db      *sql.DB
}

func NewTestEnv(t *testing.T) *TestEnv {
	os.Setenv("ENV", string(config.Env_Test))
	conf, err := config.New()
	require.NoError(t, err)
	db, err := store.NewPostgresDb(conf)
	require.NoError(t, err)


	return &TestEnv{
		Config: conf,
		Db: db,
	}
}

//test for migrations
func (te *TestEnv) SetUpDb(t *testing.T) {
	m, err := migrate.New(
		"file://../migrations",
		te.Config.DatabaseUrl())
	require.NoError(t, err)

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		require.NoError(t, err)
	}
}

