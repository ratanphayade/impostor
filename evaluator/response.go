package evaluator

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/ratanphayade/imposter/utils"
)

const (
	ResolverGenerator = "generate"
	ResolverCopy      = "copy"

	GeneratorString = "string"
	GeneratorInt    = "int"
)

var (
	generator = map[string]utils.GeneratorFunc{
		GeneratorString: utils.GenerateString,
		GeneratorInt:    utils.GenerateInt,
	}
)

type Response struct {
	Label      string            `json:"label"`
	Body       string            `json:"body"`
	Latency    int64             `json:"latency"`
	StatusCode int               `json:"status_code"`
	Headers    map[string]string `json:"headers"`
}

func (r Response) construct(data collector) Response {
	placeholder := r.parsePlaceholders()
	resolvePlaceholder(placeholder, data)
	return r.apply(placeholder)
}

func (r Response) parsePlaceholders() map[string]string {
	result := map[string]string{}
	var re = regexp.MustCompile(`(?m)({{.+?}})`)

	for _, match := range re.FindAllString(r.Body, -1) {
		result[match] = ""
	}

	return result
}

func (r Response) apply(placeholder map[string]string) Response {
	for pattern, value := range placeholder {
		m := regexp.MustCompile(pattern)
		r.Body = m.ReplaceAllString(r.Body, value)
	}

	return r
}

func resolvePlaceholder(placeholder map[string]string, data collector) {
	for k, _ := range placeholder {
		placeholder[k] = resolveTemplate(k, data)
	}
}

func resolveTemplate(key string, data collector) string {
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

func generate(tokens []string) string {
	if len(tokens) != 2 {
		log.Println("error: invalid number of arguments in generator")
		return ""
	}

	if gen, ok := generator[tokens[0]]; ok {
		val, err := strconv.Atoi(tokens[1])
		if err != nil {
			log.Println("error: ", err)
		}

		return gen(val)
	}

	return ""
}
