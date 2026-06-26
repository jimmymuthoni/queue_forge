package apiserver

import (
	"context"
	"net/http"

	"github.com/jimmymuthoni/queue_forge/config"
)

type ApiServer struct {
	Config *config.Config

}

func New(config *config.Config) *ApiServer {
	return &ApiServer{Config: config}
}

//this function ping to the server before starting
func (s *ApiServer) ping(w http.ResponseWriter, r *http.Request){
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))
}


//this function stsrts the server
func (s *ApiServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", s.ping)
	server := &http.Server{
		Addr: ":5000",
		Handler: mux,
	}
	return server.ListenAndServe()
}
