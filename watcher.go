package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"path/filepath"
	"time"
)

func watchDir(dir string, process func(file os.FileInfo, timestmap time.Time)) {
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
	watcher.Watch(dir)
}

func watchDirs(dir string, process func(dir string, file os.FileInfo, timestmap time.Time)) {
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					fmt.Println("New File:", ev.Name)
					f, _ := os.Stat(ev.Name)
					process(path.Base(path.Dir(ev.Name)), f, time.Now().Local())
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	watcher.Watch(dir)

	// Also watch insides all directories in `dir`
	d, _ := os.Open(dir)
	defer d.Close()
	files, _ := d.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			watcher.Watch(dir + string(filepath.Separator) + file.Name())
		}
	}
}
