package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"os"
	"path"
	"path/filepath"
	"time"
)

func watchDir(dir string) {
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				fmt.Println("event:", ev)
				if ev.IsCreate() {
					fmt.Println("New File:", ev.Name)
					if ValidFileType[filepath.Ext(ev.Name)] {
						ProcessFile(ev.Name, time.Now().Local())
					}
				}
				if ev.IsDelete() || ev.IsRename() {
					fmt.Println("Deleted File:", ev.Name)
					db.Where(Media{Path: ev.Name, Filename: path.Base(ev.Name)}).Delete(Media{})
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	watcher.Watch(dir)

	// Also watch inside directories
	d, _ := os.Open(dir)
	defer d.Close()
	files, _ := d.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			watcher.Watch(dir + string(filepath.Separator) + file.Name())
		}
	}
}
