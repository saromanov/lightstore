// Package watcher using for changes on the folder with config
// for applying changes without restart
package store

import (
	"github.com/fsnotify/fsnotify"
)

type watcher struct {
	watch *fsnotify.Watcher
}

func newWatcher() (*watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	return &watcher{
		watch: w,
	}, nil
}

// Close provides ending work of file watcher
func (w *watcher) Close() {
	watcher.Close()
}
