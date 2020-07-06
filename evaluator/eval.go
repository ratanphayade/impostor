package evaluator

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type Evaluator struct {
	Response Response `json:"response"`
	Rules    []Rule   `json:"rules"`
}

type Response struct {
	Label      string            `json:"label"`
	Format     string            `json:"format"`
	Latency    time.Duration     `json:"latency"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

type collector struct {
	params    map[string]string
	resources map[string]string
	headers   map[string]string
	body      map[string]interface{}
}

func Evaluate(w http.ResponseWriter, r *http.Request, evals []Evaluator) {
	data := collectRequestDetails(r)

	for _, eval := range evals {
		if eval.match(data) {
			eval.constructResponse(w)
		}
	}
}

func (e Evaluator) match(d collector) bool {
	for _, rule := range e.Rules {
		if !rule.match(d) {
			break
		}
	}
	return false
}

func (e Evaluator) constructResponse(w http.ResponseWriter) {

}

func collectRequestDetails(r *http.Request) collector {
	return collector{
		params:    collectParams(r),
		resources: collectResources(r),
		headers:   collectHeaders(r),
		body:      collectBody(r),
	}
}

func collectParams(r *http.Request) map[string]string {
	data := map[string]string{}

	for key, val := range r.URL.Query() {
		data[key] = val[0]
	}

	return data
}

func collectHeaders(r *http.Request) map[string]string {
	data := map[string]string{}

	for key, val := range r.Header {
		data[key] = val[0]
	}

	return data
}

func collectResources(r *http.Request) map[string]string {
	data := map[string]string{}

	return data
}

func collectBody(r *http.Request) map[string]interface{} {
	data := map[string]interface{}{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Print(err)
	}

	return data
}
