package main

import (
	"fmt"
	"net/http"
)

type server struct {
	host string
	port int
}

func newServer(host string, port int, cfg config) *server {
	return &server{
		host: host,
	}
}

func (s *server) run() *server {
	mux := http.NewServeMux()

	s.attachHandlers(mux)

	httpServer := http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: mux,
	}

	return s
}

func (s *server) attachHandlers(mux *http.ServeMux, routes []route) {

}
