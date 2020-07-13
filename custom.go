package main

import "log"

type callerFunc func(collector, ...string) string

const (
	CustomCustom = "custom"
)

var (
	call = map[string]callerFunc{
		CustomCustom: callCustom,
	}
)

func customCall(data collector, tokens []string) string {
	if len(tokens) < 1 {
		log.Println("error: invalid number of arguments in call")
		return ""
	}

	if caller, ok := call[tokens[0]]; ok {
		return caller(data, tokens[1:]...)
	}

	return ""
}

func callCustom(data collector, tokens ...string) string {
	return "this is a custom field"
}
