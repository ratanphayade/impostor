package config

import (
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

// Config for running the Mock server
// it will contain configs for multiple application
type Config struct {
	Apps map[string]App
}

type App struct {
	Host string
	Port int
}

// Mock Config

type MockConfig struct {
	Routes []Route
}

type Route struct {
	Method    string
	Endpoint  string
	Evaluator []Evaluator
}

type Evaluator struct {
	Response Response
	Rules    []Rule
}

type Rule struct {
	Target   string
	Modifier string
	Value    string
}

type Response struct {
	Label      string
	Format     string
	Latency    time.Duration
	StatusCode int
	Headers    map[string]string
}

var (
	Conf Config
	Mock MockConfig
)

func LoadConfig(path string, host string, port int) {
	if path == "" {
		return
	}

	if _, err := toml.DecodeFile(path, &Conf); err != nil {
		log.Fatal("failed to load Config : ", err)
	}

	Conf.Apps["default"] = App{
		Host: host,
		Port: port,
	}
}

func LoadMockConfig(path string) {
	if path == "" {
		log.Fatal("Mock Config path is not loaded")
	}

	// TODO: fix this. to load Config form below format
	//  - Mock
	//    - settlements
	//        - API-1 Config - files
	//        - API-2 Config - files
	//   also add watcher on particular dir
	if _, err := toml.DecodeFile(path, &Mock); err != nil {
		log.Fatal("failed to load Mock Config : ", err)
	}
}
