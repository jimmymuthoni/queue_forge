package apiserver_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jimmymuthoni/queue_forge/apiserver"
	"github.com/jimmymuthoni/queue_forge/config"
	"github.com/stretchr/testify/require"
)

func TestManager(t *testing.T) {
	conf, err := config.New()
	require.NoError(t, err)

	jwtManager := apiserver.NewJwtManager(conf)

	userID := uuid.New()

	tokenPair, err := jwtManager.GenerateTokenPair(userID)
	require.NoError(t, err)

	require.True(t, jwtManager.IsAccessToken(tokenPair.AccessToken))
	require.False(t, jwtManager.IsAccessToken(tokenPair.RefreshToken))

	// Access token
	accessClaims := tokenPair.AccessToken.Claims.(apiserver.CustomClaims)

	require.Equal(t, userID.String(), accessClaims.Subject)
	require.Equal(
		t,
		"http://"+conf.ApiServerHost+":"+conf.ApiServerPort,
		accessClaims.Issuer,
	)

	// Refresh token
	refreshClaims := tokenPair.RefreshToken.Claims.(apiserver.CustomClaims)

	require.Equal(t, userID.String(), refreshClaims.Subject)
	require.Equal(
		t,
		"http://"+conf.ApiServerHost+":"+conf.ApiServerPort,
		refreshClaims.Issuer,
	)
}