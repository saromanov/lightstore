// Package watcher using for changes on the folder with config
// for applying changes without restart
package store

import (
	"fmt"
	"log"

	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	watch *fsnotify.Watcher
}

func newWatcher(path string) (*watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = w.Add(path)
	if err != nil {
		return nil, err
	}

	return &watcher{
		watch: w,
	}, nil
}

// Do provides starting of the worker
func (w *watcher) Do() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-w.watch.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
			case err, ok := <-w.watch.Errors:
				if !ok {
					return
				}
				fmt.Println("Watcher error: ", err)
			}
		}
	}()
	<-done
}

// Close provides ending work of file watcher
func (w *watcher) Close() {
	w.watch.Close()
}
