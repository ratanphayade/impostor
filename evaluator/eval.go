package evaluator

import (
	"net/http"
	"time"
)

type Evaluator struct {
	Response Response `json:"response"`
	Rules    []Rule   `json:"rules"`
}

type Rule struct {
	Target   string `json:"target"`
	Modifier string `json:"modifier"`
	Value    string `json:"value"`
}

type Response struct {
	Label      string            `json:"label"`
	Format     string            `json:"format"`
	Latency    time.Duration     `json:"latency"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

func Evaluate(w http.ResponseWriter, r *http.Request, eval []Evaluator) {

}
