package server

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"log"
	"net/http"
	"time"

	"github.com/ratanphayade/imposter/evaluator"

	"github.com/gorilla/mux"
)

type App struct {
	Host string
	Port int
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

func (s *server) AttachHandlers(routes []Route) *server {
	for _, r := range routes {
		s.mux.HandleFunc(r.Endpoint, handler(r.Evaluators)).Methods(r.Method)
	}
	return s
}

func (s *server) InitializeWatcher(path string, fn func(string)) *server {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	if err := w.AddRecursive(path); err != nil {
		log.Printf("%s: error trying to watch change on %s directory", err, path)
	}

	go func() {
		if err := w.Start(time.Millisecond * 1000); err != nil {
			log.Println(err)
		}
	}()

	readEventsFromWatcher(w, path, fn)

	return s
}

func readEventsFromWatcher(w *watcher.Watcher, path string,  fn func(path string)) {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				log.Println("Modified file:", evt.Name())
				fn(path)
			case err := <-w.Error:
				log.Println("error checking file change:", err)
				fn(path)
			case <-w.Closed:
				return
			}
		}
	}()
}

func handler(evals []evaluator.Evaluator) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		evaluator.Evaluate(w, r, evals)
	}
}
