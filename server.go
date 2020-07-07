package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var (
	defaultCORSMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodOptions,
		http.MethodDelete,
		http.MethodPatch,
		http.MethodTrace,
		http.MethodConnect}

	defaultCORSHeaders = []string{"X-Requested-With", "Content-Type", "Authorization"}

	defaultCORSExposedHeaders = []string{"Cache-Control", "Content-Language", "Content-Type", "Expires", "Last-Modified", "Pragma"}
)

type App struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	MockPath string `toml:"mock_path"`
}

// Mock Config
type MockConfig struct {
	Routes      []Route
	NotFound    Response
	CORSOptions CORSOption
}

type CORSOption struct {
	Methods          []string `json:"methods"`
	Headers          []string `json:"headers"`
	ExposedHeaders   []string `json:"exposed_headers"`
	Origins          []string `json:"origins"`
	AllowCredentials bool     `json:"allow_credentials"`
}

type Route struct {
	Method     string      `json:"method"`
	Endpoint   string      `json:"endpoint"`
	Evaluators []Evaluator `json:"evaluators"`
}

type server struct {
	app        App
	httpServer *http.Server
	mux        *mux.Router
	mock       MockConfig
}

func NewServer(cfg map[string]App, application string, mock MockConfig) *server {
	var (
		listerConf App
		ok         bool
	)

	if listerConf, ok = cfg[application]; !ok {
		listerConf = cfg["default"]
	}

	return &server{
		app:  listerConf,
		mux:  mux.NewRouter(),
		mock: mock,
	}
}

func (s *server) AttachHandlers() *server {
	for _, r := range s.mock.Routes {
		s.mux.HandleFunc(r.Endpoint, handler(r.Evaluators, s.mock.NotFound)).
			Methods(r.Method)
	}

	s.mux.NotFoundHandler = s.mux.NewRoute().
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, s.mock.NotFound.Body, http.StatusNotFound)
		}).GetHandler()

	return s
}

func (s *server) refresh(route Route) {
	// todo: refresh the route list available in the changed file
}

func (s *server) Run() *server {
	s.httpServer = &http.Server{
		Handler: handlers.CORS(collectCORSOptions(s.mock.CORSOptions)...)(s.mux),
		Addr:    fmt.Sprintf("%s:%d", s.app.Host, s.app.Port),
	}

	if err := s.httpServer.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

	return s
}

func handler(evals []Evaluator, notFound Response) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res := Evaluate(r, evals, notFound)

		writeResponse(w, res)
	}
}

func writeResponse(w http.ResponseWriter, res Response) {
	for k, v := range res.Headers {
		w.Header().Set(k, v)
	}

	if res.StatusCode >= http.StatusBadRequest && res.StatusCode <= http.StatusNotExtended {
		http.Error(w, res.Body, res.StatusCode)
		return
	}

	w.WriteHeader(res.StatusCode)
	_, _ = w.Write([]byte(res.Body))

	time.Sleep(time.Duration(res.Latency) * time.Millisecond)
}

func collectCORSOptions(cors CORSOption) []handlers.CORSOption {
	var h []handlers.CORSOption

	if len(cors.Methods) > 0 {
		h = append(h, handlers.AllowedMethods(cors.Methods))
	} else {
		h = append(h, handlers.AllowedMethods(defaultCORSMethods))
	}

	if len(cors.Origins) > 0 {
		h = append(h, handlers.AllowedOrigins(cors.Origins))
	}

	if len(cors.Headers) > 0 {
		h = append(h, handlers.AllowedHeaders(cors.Headers))
	} else {
		h = append(h, handlers.AllowedHeaders(defaultCORSHeaders))
	}

	if len(cors.ExposedHeaders) > 0 {
		h = append(h, handlers.ExposedHeaders(cors.ExposedHeaders))
	} else {
		h = append(h, handlers.ExposedHeaders(defaultCORSExposedHeaders))
	}

	if cors.AllowCredentials {
		h = append(h, handlers.AllowCredentials())
	}

	return h
}
