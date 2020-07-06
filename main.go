package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/ratanphayade/imposter/evaluator"

	"github.com/BurntSushi/toml"
	"github.com/ratanphayade/imposter/server"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 9000
	defaultRequestsPath = ""
	defaultConfigFile   = ""
	defaultApplication  = ""
	defaultWatch        = false

	NotFoundResponseFile = "404.json"
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

var (
	Conf Config
	Mock server.MockConfig
)

func init() {
	host = flag.String("host", defaultHost, "if you run your server on a different host")
	port = flag.Int("port", defaultPort, "port to run the server")
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in application mock config and reload")
	mockPath = flag.String("mock", defaultRequestsPath, "directory where your mock configs are saved")
	application = flag.String("application", defaultApplication, "name of the application which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	LoadConfig(*configFilePath, *host, *port)
	LoadMockConfig(*mockPath)
}

func main() {
	server.NewServer(Conf.Apps, *application).
		AttachHandlers(Mock).
		Run()
}

func LoadConfig(path string, host string, port int) {
	Conf.Apps = make(map[string]server.App)

	Conf.Apps["default"] = server.App{
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
	if path == "" {
		if cnf, ok := Conf.Apps[*application]; ok {
			path = cnf.MockPath
		} else {
			log.Fatal("invalid mock config path")
		}
	}

	path = path + "/" + *application

	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal("failed to load open Mock Directory : ", err)
	}

	var routes []server.Route

	for _, v := range files {
		filePath := path + "/" + v.Name()

		if v.Name() == NotFoundResponseFile {
			var dest evaluator.Response
			readRequestMockConfig(filePath, &dest)
			Mock.NotFound = dest
		} else {
			var route server.Route
			readRequestMockConfig(filePath, &route)
			routes = append(routes, route)
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
