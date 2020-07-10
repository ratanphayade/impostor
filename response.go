package main

import (
	"log"
	"regexp"
	"strings"
)

const (
	ResolverGenerator = "generate"
	ResolverCopy      = "copy"
	ResolverCall      = "call"
)

type Response struct {
	Label      string            `json:"label"`
	Body       string            `json:"body"`
	Latency    int64             `json:"latency"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

func (r Response) construct(data collector) Response {
	log.Println("constructing response for label: ", r.Label)

	placeholders := regexp.MustCompile(`(?m)({{.+?}})`).FindAllString(r.Body, -1)

	for _, placeholder := range placeholders {
		result := resolvePlaceholder(placeholder, data)
		r.Body = strings.Replace(r.Body, placeholder, result, 1)
	}

	return r
}

func resolvePlaceholder(key string, data collector) string {
	tokens := tokenize(key)

	return resolve(tokens[0], tokens[1:], data)
}

func tokenize(str string) []string {
	tokens := regexp.MustCompile(`(?m)(\s*\S+\s*)?`).FindAllString(str, -1)
	tokens = tokens[1 : len(tokens)-1]

	for i, v := range tokens {
		tokens[i] = strings.TrimSpace(v)
	}

	return tokens
}

func resolve(key string, tokens []string, data collector) string {
	switch key {
	case ResolverGenerator:
		return generate(tokens)

	case ResolverCopy:
		return copyFrom(tokens, data)

	case ResolverCall:
		return customCall(data, tokens)
	}

	return ""
}

func copyFrom(tokens []string, data collector) string {
	if len(tokens) != 2 {
		log.Println("error: invalid number of arguments in copy")
		return ""
	}

	return data.get(tokens[0], tokens[1])
}
