package apiserver

import (
	"log/slog"
	"net/http"
)

//this function logs all the requests
func NewLoggerMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler{
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				logger.Info("http request","method", r.Method, "path", r.URL.Path)
				next.ServeHTTP(w, r)
			})
		}
}