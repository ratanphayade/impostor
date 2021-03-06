package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"

	"github.com/BurntSushi/toml"
)

var (
	// this is status code order based on which we will be ordering the
	StatusCodeOrder = map[int]int{
		http.StatusInternalServerError: 7,
		http.StatusUnauthorized:        6,
		http.StatusBadRequest:          5,
		http.StatusForbidden:           4,
		http.StatusNotFound:            3,
		http.StatusCreated:             2,
		http.StatusOK:                  1,
	}
)

func LoadConfig(path string, host string, port int) {
	Conf.Apps = make(map[string]App)

	Conf.Apps["default"] = App{
		Host:     host,
		Port:     port,
		MockPath: "test",
	}

	if path != "" {
		if _, err := toml.DecodeFile(path, &Conf); err != nil {
			log.Fatal("failed to load Config : ", err)
		}
	}
}

func LoadMockConfig(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("failed to load open Mock Directory : ", err)
	}

	var routes []Route

	for _, v := range files {
		if filepath.Ext(v.Name()) != ".json" {
			log.Println("error: invalid file format: ", v.Name())
			continue
		}

		filePath := path + string(os.PathSeparator) + v.Name()

		switch v.Name() {
		case CORSFile:
			var cors CORSOption
			readRequestMockConfig(filePath, &cors)
			Mock.CORSOptions = cors

		case NotFoundResponseFile:
			var notFound Response
			readRequestMockConfig(filePath, &notFound)
			Mock.NotFound = notFound

		default:
			var route Route
			readRequestMockConfig(filePath, &route)
			routes = append(routes, route)

			evals := route.Evaluators
			sort.Slice(evals, func(i, j int) bool {
				return StatusCodeOrder[evals[i].Response.StatusCode] >
					StatusCodeOrder[evals[j].Response.StatusCode]
			})

			route.Evaluators = evals
		}
	}

	if len(routes) == 0 {
		log.Println("failed to load route configurations")
	}

	Mock.Routes = routes
}

func readRequestMockConfig(filePath string, dest interface{}) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Println("failed to load Mock Config : ", err)
		return
	}

	if err := json.Unmarshal(file, dest); err != nil {
		log.Println("failed to un marshal Mock Config : ", err)
		return
	}
}
