package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/nqd/flat"

	"github.com/gorilla/mux"
)

type Evaluator struct {
	Response Response `json:"response"`
	Rules    []Rule   `json:"rules"`
}

func Evaluate(r *http.Request, evals []Evaluator, notFound Response) Response {
	var (
		hasDefault  bool
		defaultEval Evaluator
		data        = collectRequestDetails(r)
	)

	for _, eval := range evals {
		if def, match := eval.match(data); match {
			if def {
				hasDefault = true
				defaultEval = eval
			} else {
				return eval.Response.construct(data)
			}
		}
	}

	if hasDefault {
		return defaultEval.Response.construct(data)
	}

	return notFound
}

func (e Evaluator) match(d collector) (bool, bool) {
	if len(e.Rules) == 0 {
		return true, true
	}

	for _, rule := range e.Rules {
		if !rule.match(d) {
			return false, false
		}
	}
	return false, true
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
	request := map[string]interface{}{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Println(err)
	}

	data, err := flat.Flatten(request, nil)
	if err != nil {
		log.Println(err)
	}

	return data
}

func get(val map[string]interface{}, key string) string {
	var result interface{}

	if v, ok := val[key]; ok {
		result = v
	} else {
		log.Println("key not found ", key)
		return ""
	}

	return fmt.Sprintf("%v", result)
}
