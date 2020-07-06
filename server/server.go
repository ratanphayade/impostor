package server

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ratanphayade/imposter/evaluator"

	"github.com/gorilla/mux"
)

type App struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	MockPath string `toml:"mock_path"`
}

// Mock Config
type MockConfig struct {
	Routes   []Route
	NotFound evaluator.Response
}

type Route struct {
	Method     string                `json:"method"`
	Endpoint   string                `json:"endpoint"`
	Evaluators []evaluator.Evaluator `json:"evaluators"`
}

type server struct {
	app App
	mux *mux.Router
}

func NewServer(cfg map[string]App, application string) *server {
	var (
		listerConf App
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

func (s *server) AttachHandlers(mc MockConfig) *server {
	for _, r := range mc.Routes {
		s.mux.HandleFunc(r.Endpoint, handler(r.Evaluators, mc.NotFound)).
			Methods(r.Method)
	}

	s.mux.NotFoundHandler = s.mux.NewRoute().
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, mc.NotFound.Body, http.StatusNotFound)
		}).GetHandler()

	return s
}

func handler(evals []evaluator.Evaluator, notFound evaluator.Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := evaluator.Evaluate(r, evals, notFound)

		writeResponse(w, res)
	}
}

func writeResponse(w http.ResponseWriter, res evaluator.Response) {
	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}

	if res.StatusCode == http.StatusNotFound {
		http.Error(w, res.Body, http.StatusNotFound)
		return
	} else {
		_, _ = w.Write([]byte(res.Body))
	}

	w.WriteHeader(res.StatusCode)

	time.Sleep(time.Duration(res.Latency) * time.Millisecond)
}
