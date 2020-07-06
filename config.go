package main

import "time"

// config for running the mock server
// it will contain configs for multiple application
type config struct {
	apps map[string]app
}

type app struct {
	host string
	port int
}

// mock config

type mockConfig struct {
	routes []route
}

type route struct {
	method    string
	endpoint  string
	evaluator []evaluator
}

type evaluator struct {
	response response
	rules    []rule
}

type rule struct {
	target   string
	modifier string
	value    string
}

type response struct {
	label      string
	format     string
	latency    time.Duration
	statusCode int
	headers    map[string]string
}
