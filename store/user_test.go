package store_test

import (
	"context"
	"testing"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jimmymuthoni/queue_forge/fixtures"
	"github.com/stretchr/testify/require"
	"github.com/jimmymuthoni/queue_forge/store"
)

func TestUserStore(t *testing.T){

	env := fixtures.NewTestEnv(t)
	env.SetUpDb(t)
	cleanup := env.SetUpDb(t)
	t.Cleanup(func() {
		cleanup(t)
	})

	now := time.Now()
	ctx := context.Background()
	userStore := store.NewUserStore(env.Db)
	user, err := userStore.CreateUser(context.Background(), "jim@testingmail.com", "testingpassword")
	require.NoError(t, err)

	require.Equal(t, "jim@testingmail.com", user.Email)
	require.NoError(t, user.ComparePassword("testingpassword"))
	require.Less(t, now.UnixNano(), user.CreatedAt.UnixNano())

	//teting fnding user byid
	user2, err := userStore.FindUserById(ctx, user.Id)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

	//testing finding user by email
	user2, err = userStore.FindUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.Equal(t, user.Email, user2.Email)
	require.Equal(t, user.Id, user2.Id)
	require.Equal(t, user.HashedPasswordBase64, user2.HashedPasswordBase64)
	require.Equal(t, user.CreatedAt.UnixNano(), user2.CreatedAt.UnixNano())

}