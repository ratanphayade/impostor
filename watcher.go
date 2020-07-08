package main

import (
	"log"
	"time"

	"github.com/radovskyb/watcher"
)

func initializeWatcher(path string, s *server) {
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

	readEventsFromWatcher(w, path, s)
}

func readEventsFromWatcher(w *watcher.Watcher, path string, s *server) {
	go func() {
		for {
			select {
			case evt := <-w.Event:
				log.Println("modified file:", evt.Name())
				LoadMockConfig(path)
				s.refresh()

			case err := <-w.Error:
				log.Println("error checking file change:", err)
				LoadMockConfig(path)

			case <-w.Closed:
				return
			}
		}
	}()
}
