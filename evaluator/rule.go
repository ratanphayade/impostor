package evaluator

import (
	"log"
	"regexp"
)

const (
	TargetResource = "resource"
	TargetParams   = "params"
	TargetHeader   = "header"
	TargetBody     = "body"
)

type Rule struct {
	Target   string `json:"target"`
	Modifier string `json:"modifier"`
	Value    string `json:"value"`
	IsRegex  bool   `json:"is_regex"`
}

func (rule Rule) match(d collector) bool {
	val := rule.getValue(d)

	if rule.IsRegex {
		exp, err := regexp.Compile(rule.Value)
		if err != nil {
			log.Fatal(err)
		}
		return exp.Match([]byte(rule.Value))
	}

	return val == rule.Value
}

func (rule Rule) getValue(d collector) string {
	switch rule.Target {
	case TargetParams:
		return d.getFromParam(rule.Modifier)

	case TargetHeader:
		return d.getFromHeader(rule.Modifier)

	case TargetBody:
		return d.getFromBody(rule.Modifier)

	case TargetResource:
		return d.getFromResource(rule.Modifier)
	}

	log.Println("no matching target found")

	return ""
}
