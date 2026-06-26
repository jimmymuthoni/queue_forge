package apiserver

import "github.com/jimmymuthoni/queue_forge/config"

type ApiServer struct {
	Config *config.Config

}


func New(config *config.Config) *ApiServer {
	return &ApiServer{Config: config}
}

