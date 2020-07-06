package evaluator

import (
	"net/http"
	"time"
)

type Evaluator struct {
	Response Response
	Rules    []Rule
}

type Rule struct {
	Target   string
	Modifier string
	Value    string
}

type Response struct {
	Label      string
	Format     string
	Latency    time.Duration
	StatusCode int
	Headers    map[string]string
}

func Evaluate(w http.ResponseWriter, r *http.Request, eval []Evaluator) {

}
