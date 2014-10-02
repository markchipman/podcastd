package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"time"
)

type Movie struct {
	Title     string
	Year      int
	Filename  string
	Size      int
	Timestamp time.Time
}

type TVShow struct {
	Title     string
	Season    int
	Episode   int
	Aired     time.Time
	Filename  string
	Size      int
	Timestamp time.Time
}

func (t TVShow) TableName() string {
	return "tvshows"
}

func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", dir+"btpodcast.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	//db.LogMode(true)
	db.AutoMigrate(&Movie{}, &TVShow{})
	return db
}

var db = initDB()

func updateDB() {
	timestamp := time.Now()
	d, _ := os.Open(dir + "movies")
	defer d.Close()
	files, _ := d.Readdir(-1)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".m4v" {
			movie := Movie{}
			db.FirstOrCreate(&movie, Movie{
				Filename:  file.Name(),
				Size:      int(file.Size()) / 1024 / 1024,
				Timestamp: timestamp,
			})
		}
	}
	d, _ = os.Open(dir + "tvshows")
	defer d.Close()
	files, _ = d.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			e, _ := os.Open(dir + "tvshows" + string(filepath.Separator) + file.Name())
			defer e.Close()
			episodes, _ := e.Readdir(-1)
			for _, episode := range episodes {
				if filepath.Ext(episode.Name()) == ".m4v" {
					show := TVShow{}
					db.FirstOrCreate(&show, TVShow{
						Title:     file.Name(),
						Filename:  episode.Name(),
						Size:      int(episode.Size()) / 1024 / 1024,
						Timestamp: timestamp,
					})
				}
			}
		}
	}
}
