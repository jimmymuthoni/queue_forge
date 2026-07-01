package apiserver

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"
	"fmt"

	"github.com/jimmymuthoni/queue_forge/config"
	"github.com/jimmymuthoni/queue_forge/store"
)

type ApiServer struct {
	config *config.Config
	logger *slog.Logger
	store  *store.Store
}

func New(config *config.Config, logger *slog.Logger, store *store.Store) *ApiServer {
	return &ApiServer{
		config: config,
		logger: logger,
		store:  store,
	}
}

// this function ping to the server before starting
func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}

// this function starts the server and logs all the request
func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", s.ping)
	mux.HandleFunc("POST /auth/signup", s.signupHandler)

	// middleware := NewLoggerMiddleware(s.logger)

	server := &http.Server{
    Addr:    net.JoinHostPort(s.config.ApiServerHost, s.config.ApiServerPort),
    Handler: mux,
}

	//goroutine to handle start logic of the server
	go func() {
		slog.Info("Server running", "port", s.config.ApiServerPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Error("Api server failed to listen and serve", "error", err)
		}
	}()

	//go rountine to handle the shutdown logic of the sever
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		<-ctx.Done()

		shutDownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutDownCtx); err != nil {
			s.logger.Error("APiserver faile to shutdown", "error", err)
		}
	}()
	wg.Wait()

	return nil
}
