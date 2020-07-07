package evaluator

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Evaluator struct {
	Response Response `json:"response"`
	Rules    []Rule   `json:"rules"`
}

type values map[string]string

func (v values) get(key string) string {
	if val, ok := v[key]; ok {
		return val
	}

	return ""
}

type collector struct {
	params    values
	resources values
	headers   values
	body      map[string]interface{}
}

func (c collector) getFromParam(key string) string {
	return c.params.get(key)
}

func (c collector) getFromResource(key string) string {
	return c.resources.get(key)
}

func (c collector) getFromHeader(key string) string {
	return c.headers.get(key)
}

func (c collector) getFromBody(key string) string {
	return get(c.body, key)
}

func (c collector) get(target string, modifier string) string {
	switch target {
	case TargetParams:
		return c.getFromParam(modifier)

	case TargetHeader:
		return c.getFromHeader(modifier)

	case TargetBody:
		return c.getFromBody(modifier)

	case TargetResource:
		return c.getFromResource(modifier)
	}

	log.Println("no matching target found")

	return ""
}

func Evaluate(r *http.Request, evals []Evaluator, notFound Response) Response {
	data := collectRequestDetails(r)

	for _, eval := range evals {
		if eval.match(data) {
			return eval.Response.construct(data)
		}
	}

	return notFound
}

func (e Evaluator) match(d collector) bool {
	for _, rule := range e.Rules {
		if !rule.match(d) {
			return false
		}
	}
	return true
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

	for key, val := range mux.Vars(r) {
		data[key] = val
	}

	return data
}

func collectBody(r *http.Request) map[string]interface{} {
	data := map[string]interface{}{}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		log.Println(err)
	}

	return data
}

func get(val map[string]interface{}, key string) string {
	var (
		result interface{}
		keys   = strings.Split(key, ".")
	)

	if v, ok := val[keys[0]]; ok {
		result = v
	} else {
		return ""
	}

	if len(keys) == 1 {
		return fmt.Sprintf("%v", result)
	}

	if subVal, ok := result.(map[string]interface{}); ok {
		return get(subVal, strings.Join(keys[1:], "."))
	}

	return ""
}
