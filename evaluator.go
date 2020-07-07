package main

import (
	"net/http"
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
