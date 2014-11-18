package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"os"
	"time"
)

func watchDir(path string, process func(file os.FileInfo, timestmap time.Time)) {
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					fmt.Println("New File:", ev.Name)
					f, _ := os.Stat(ev.Name)
					process(f, time.Now().Local())
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	watcher.Watch(path)
}
