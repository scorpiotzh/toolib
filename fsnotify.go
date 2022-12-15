package toolib

import (
	"github.com/fsnotify/fsnotify"
	"log"
	"path/filepath"
)

func AddFileWatcher(filePath string, handle func()) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	configFile := filepath.Clean(filePath)
	configDir, _ := filepath.Split(configFile)
	log.Println("AddFileWatcher:", configDir)
	err = watcher.Add(filePath)
	if err != nil {
		return watcher, err
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				log.Println("event:", event)
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Chmod == fsnotify.Chmod {
					log.Println("modified file:", event.Name)
					handle()
				} else if event.Op&fsnotify.Rename == fsnotify.Rename {
					log.Println("modified file:", event.Name)
					_ = watcher.Remove(filePath)
					_ = watcher.Add(filePath)
				}
			case err, ok := <-watcher.Errors:
				log.Println("error:", err)
				if !ok {
					return
				}

			}
		}
	}()
	return watcher, nil
}
