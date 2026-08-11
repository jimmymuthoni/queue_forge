package store_test

import (
	"context"
	"testing"

	"github.com/jimmymuthoni/queue_forge/apiserver"
	"github.com/jimmymuthoni/queue_forge/fixtures"
	"github.com/jimmymuthoni/queue_forge/store"
	"github.com/stretchr/testify/require"
)

func TestRefreshTokenStore(t *testing.T){
	env := fixtures.NewTestEnv(t)
	cleanup := env.SetUpDb(t)
	t.Cleanup(func ()  {
		cleanup(t)
	})

	ctx := context.Background()
	refreshTokenStore := store.NewRefreshTokenStore(env.Db)
	userStore := store.NewUserStore(env.Db)
	user, err := userStore.CreateUser(ctx, "test@gmail.com", "test")
	require.NoError(t, err)

	jwtManager := apiserver.NewJwtManager(env.Config)
	tokenPair, err := jwtManager.GenerateTokenPair(user.Id)
	require.NoError(t, err)


	refreshTokenRecord, err := refreshTokenStore.Create(ctx,user.Id, tokenPair.RefreshToken)
	require.NoError(t, err)
	require.Equal(t, user.Id, refreshTokenRecord.UserId)
	expectedExpiration, err := tokenPair.RefreshToken.Claims.GetExpirationTime()
	require.Equal(t, expectedExpiration.Time.UnixMilli(),refreshTokenRecord.ExpiresAt.UnixMilli())

	refreshTokenRecord2, err := refreshTokenStore.ByPrimaryKey(ctx, user.Id, tokenPair.RefreshToken)
	require.Equal(t, refreshTokenRecord.UserId, refreshTokenRecord2.UserId)
	require.Equal(t, refreshTokenRecord.HashedToken, refreshTokenRecord2.HashedToken)
	require.Equal(t, refreshTokenRecord.CreatedAt, refreshTokenRecord2.CreatedAt)
	require.Equal(t, refreshTokenRecord.ExpiresAt, refreshTokenRecord2.ExpiresAt)

	result, err := refreshTokenStore.DeleteUserTokens(ctx ,user.Id)
	require.NoError(t, err)
	rowsAffected, err := result.RowsAffected()
	require.NoError(t, err)
	require.Equal(t, int64(1), rowsAffected)


}

