package main

import (
	"flag"

	"github.com/ratanphayade/impostor/server"
)

const (
	defaultHost         = "localhost"
	defaultPort         = 9000
	defaultRequestsPath = ""
	defaultConfigFile   = ""
	defaultApplication  = ""
	defaultWatch        = false

	NotFoundResponseFile = "404.json"
	CORSFile             = "cors.json"
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
	Conf        Config
	appMockPath string
	Mock        server.MockConfig
)

func init() {
	host = flag.String("host", defaultHost, "if you run your server on a different host")
	port = flag.Int("port", defaultPort, "port to run the server")
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in application mock config and reload")
	mockPath = flag.String("mock", defaultRequestsPath, "directory where your mock configs are saved")
	application = flag.String("application", defaultApplication, "name of the application which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	appMockPath = *mockPath + "/" + *application
	LoadConfig(*configFilePath, *host, *port)
	LoadMockConfig(appMockPath)
}

func main() {
	if *watch {
		initializeWatcher(appMockPath)
	}

	server.NewServer(Conf.Apps, *application, Mock).
		AttachHandlers().
		Run()
}
