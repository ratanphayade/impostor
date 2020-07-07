package main

import (
	"flag"

	"github.com/ratanphayade/impostor/server"
)

const (
	defaultHost        = "localhost"
	defaultPort        = 9000
	defaultMockPath    = "test"
	defaultConfigFile  = "config.toml"
	defaultApplication = "app"
	defaultWatch       = false

	NotFoundResponseFile = "404.json"
	CORSFile             = "cors.json"
)

var (
	host           *string
	port           *int
	mockPath       *string
	app            *string
	configFilePath *string
	watch          *bool
)

// Config for running the Mock server
// it will contain configs for multiple app
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
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in app mock config and reload")
	mockPath = flag.String("mock", defaultMockPath, "directory where your mock configs are saved")
	app = flag.String("app", defaultApplication, "name of the app which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	appMockPath = *mockPath + "/" + *app
	LoadConfig(*configFilePath, *host, *port)
	LoadMockConfig(appMockPath)
}

func main() {
	if *watch {
		initializeWatcher(appMockPath)
	}

	server.NewServer(Conf.Apps, *app, Mock).
		AttachHandlers().
		Run()
}
