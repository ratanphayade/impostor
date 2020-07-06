package main

import (
	"flag"

	"github.com/ratanphayade/imposter/server"

	"github.com/ratanphayade/imposter/config"
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

func init() {
	host = flag.String("host", defaultHost, "if you run your server on a different host")
	port = flag.Int("port", defaultPort, "port to run the server")
	watch = flag.Bool("watch", defaultWatch, "if true, then watch for any change in application mock config and reload")
	mockPath = flag.String("mocks", defaultRequestsPath, "directory where your mock configs are saved")
	application = flag.String("application", defaultApplication, "name of the application which has to be mocked")
	configFilePath = flag.String("config", defaultConfigFile, "path with configuration file")

	flag.Parse()

	config.LoadConfig(*configFilePath, *host, *port)
	config.LoadMockConfig(*mockPath)
}

func main() {
	server.NewServer(config.Conf.Apps, *application).
		AttachHandlers(config.Mock.Routes).
		Run()
}
