package main

import (
	"flag"
	"log"

	"github.com/BurntSushi/toml"
	_ "github.com/radovskyb/watcher"
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

	conf config
	mock mockConfig
)

func init() {
	host = flag.String("host", defaultHost, "if you run your server on a different host")
	port = flag.Int("port", defaultPort, "port to run the server")
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in application mock config and reload")
	mockPath = flag.String("mocks", defaultRequestsPath, "directory where your mock configs are saved")
	application = flag.String("application", defaultApplication, "name of the application which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	loadConfig()
	loadMockConfig()
}

func main() {
	srv := &server{}
	srv.run(*host, *port, conf)
}

func loadConfig() {
	if *configFilePath == "" {
		return
	}

	if _, err := toml.DecodeFile(*configFilePath, &conf); err != nil {
		log.Fatal("failed to load config : ", err)
	}
}

func loadMockConfig() {
	if *mockPath == "" {
		log.Fatal("mock config path is not loaded")
	}

	// TODO: fix this. to load config form below format
	//  - mock
	//    - settlements
	//        - API-1 config - files
	//        - API-2 config - files
	//   also add watcher on particular dir
	if _, err := toml.DecodeFile(*mockPath, &mock); err != nil {
		log.Fatal("failed to load mock config : ", err)
	}
}
