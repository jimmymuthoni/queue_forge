package fixtures

import (
	"database/sql"
	"os"
	"testing"

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


