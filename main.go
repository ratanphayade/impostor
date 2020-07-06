package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/radovskyb/watcher"
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

	initializeWatcher(appMockPath)

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

func initializeWatcher(path string) {
	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Rename, watcher.Move, watcher.Create, watcher.Write)

	if err := w.Add(path); err != nil {
		log.Printf("%s: error trying to watch change on %s directory", err, path)
	}

	go func() {
		if err := w.Start(time.Second); err != nil {
			log.Println(err)
		}
	}()

	readEventsFromWatcher(w, path)
}

func readEventsFromWatcher(w *watcher.Watcher, path string) {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				log.Println("modified file:", evt.Name())
				LoadMockConfig(path)

			case err := <-w.Error:
				log.Println("error checking file change:", err)
				LoadMockConfig(path)

			case <-w.Closed:
				return
			}
		}
	}()
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
