package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	"github.com/ratanphayade/imposter/evaluator"

	"github.com/ratanphayade/imposter/server"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 9000
	defaultRequestsPath = ""
	defaultConfigFile   = ""
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
	//if path == "" {
	//	log.Fatal("Mock Config path is not loaded")
	//}
	//
	//// TODO: fix this. to load Config form below format
	////  - Mock
	////    - settlements
	////        - API-1 Config - files
	////        - API-2 Config - files
	////   also add watcher on particular dir
	//if _, err := toml.DecodeFile(path, &Mock); err != nil {
	//	log.Fatal("failed to load Mock Config : ", err)
	//}

	Mock.Routes = []server.Route{
		{
			Method:   "GET",
			Endpoint: "/users",
			Evaluators: []evaluator.Evaluator{
				{
					Response: evaluator.Response{
						Label:      "success",
						Format:     `{"name": "Ratan Phayade"}`,
						Latency:    0,
						StatusCode: 0,
						Headers:    map[string]string{},
					},
					Rules: []evaluator.Rule{
						{
							Target:   "",
							Modifier: "",
							Value:    "",
						},
					},
				},
			},
		},
	}
}
