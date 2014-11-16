package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
)

func watchDownloads() {
	watcher, _ := fsnotify.NewWatcher()
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() {
					fmt.Println("New Download:", ev.Name)
				}
			case err := <-watcher.Error:
				fmt.Println("error:", err)
			}
		}
	}()
	watcher.Watch(config.Movies)
}
