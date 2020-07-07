package main

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
		return exp.Match([]byte(val))
	}

	return val == rule.Value
}

func (rule Rule) getValue(d collector) string {
	return d.get(rule.Target, rule.Modifier)
}
