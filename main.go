package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/ratanphayade/imposter/server"
	"io/ioutil"
	"log"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 9000
	defaultRequestsPath = "mock"
	defaultConfigFile   = "config.toml"
	defaultApplication  = ""
	defaultWatch        = false
)

var (
	host           *string
	port           *int
	mockPath       *string
	application    *string
	configFilePath *string
	watch          *bool
)

// Config for running the Mock server
// it will contain configs for multiple application
type Config struct {
	Apps map[string]server.App
}

// Mock Config
type MockConfig struct {
	Routes []server.Route
}

var (
	Conf Config
	Mock MockConfig
)

func init() {
	host = flag.String("host", defaultHost, "if you run your server on a different host")
	port = flag.Int("port", defaultPort, "port to run the server")
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in application mock config and reload")
	mockPath = flag.String("mocks", defaultRequestsPath, "directory where your mock configs are saved")
	application = flag.String("application", defaultApplication, "name of the application which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	LoadConfig(*configFilePath, *host, *port)
	LoadMockConfig(*mockPath)
}

func main() {
	server.NewServer(Conf.Apps, *application).
		AttachHandlers(Mock.Routes).
		Run()
}

func LoadConfig(path string, host string, port int) {
	Conf.Apps = make(map[string]server.App)

	Conf.Apps["default"] = server.App{
		Host: host,
		Port: port,
	}

	if path != "" {
		if _, err := toml.DecodeFile(path, &Conf); err != nil {
			log.Fatal("failed to load Config : ", err)
		}
	}
}

func LoadMockConfig(path string) {
	path = path + "/" + *application

	files , err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("failed to load open Mock Directory : ", err)
	}

	var routes []server.Route

	for _, v := range files {
		var route server.Route
		file , err := ioutil.ReadFile(path+"/"+v.Name())
		if err != nil {
			log.Fatal("failed to load Mock Config : ", err)
		}

		if err := json.Unmarshal(file, &route); err != nil {
			log.Fatal("failed to un marshal Mock Config : ", err)
		}

		routes = append(routes, route)
	}

	Mock.Routes = routes
}
