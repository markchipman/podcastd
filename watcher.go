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
					fp, fn := path.Split(ev.Name)
					db.Where(Media{Path: fp, Filename: fn}).Delete(Media{})
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() {
			watcher.Watch(path)
		}
		return nil
	})
}
