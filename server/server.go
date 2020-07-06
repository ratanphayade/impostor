package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/ratanphayade/imposter/config"

	"github.com/gorilla/mux"
)

type server struct {
	app config.App
	mux *mux.Router
}

func NewServer(cfg map[string]config.App, application string) *server {
	var (
		listerConf config.App
		ok         bool
	)

	if listerConf, ok = cfg[application]; !ok {
		listerConf = cfg["default"]
	}

	return &server{
		app: listerConf,
		mux: mux.NewRouter(),
	}
}

func (s *server) Run() *server {
	httpServer := http.Server{
		Handler: s.mux,
		Addr:    fmt.Sprintf("%s:%d", s.app.Host, s.app.Port),
	}

	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	return s
}

func (s *server) AttachHandlers(routes []config.Route) *server {
	for _, r := range routes {
		s.mux.HandleFunc(r.Endpoint, handler(r.Evaluator)).Methods(r.Method)
	}
	return s
}

func handler(eval []config.Evaluator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
