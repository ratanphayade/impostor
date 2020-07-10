package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type GeneratorFunc func(n int) string

const (
	GeneratorString = "string"
	GeneratorInt    = "int"
)

var (
	generator = map[string]GeneratorFunc{
		GeneratorString: generateString,
		GeneratorInt:    generateInt,
	}

	letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

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

	log.Println("error: specified generator not found")
	return ""
}

func generateString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateInt(n int) string {
	return fmt.Sprintf("%d", rand.Int63n(1e16))[0:n]
}
