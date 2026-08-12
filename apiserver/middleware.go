package apiserver

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/jimmymuthoni/queue_forge/store"
)

// this function logs all the requests
func NewLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info("http request", "method", r.Method, "path", r.URL.Path)
			next.ServeHTTP(w, r)
		})
	}
}

// authentication middleware
func NewAuthMiddleware(jwtManager *JwtManager, userStore *store.UserStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/auth"){
				next.ServeHTTP(w, r)
			}
			//check auth headers
			authHeader := r.Header.Get("/Authorization")
			var token string
			if parts := strings.Split(authHeader, "Bearer"); len(parts) == 2 { 
				token = parts[1]
			}
			if token == ""{
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}
			parsedToken, err := jwtManager.Parse(token)
			if err != nil {
				slog.Error("failed to parse the token", "error", err)
				w.WriteHeader(http.StatusInternalServerError)
				return 

			}

			if !jwtManager.IsAccessToken(parsedToken) {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("not an access token"))
				return 
			}

			userIsStr, err := parsedToken.Claims.GetSubject()
			if err != nil {
				slog.Error("faled to extract subject claim from token","error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}

			userId, err := uuid.Parse(userIsStr)
			if err != nil {
				slog.Error("token subject is not valid uuid", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}

			user, err := userStore.FindUserById(r.Context(), userId)
			if err != nil {
				slog.Error("faied to get user by id", "error", err)
				w.WriteHeader(http.StatusUnauthorized)
				return 
			}

			
		})
	}
}
