package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type Media interface {
	MediaTitle() string
	MediaFilename() string
	MediaType() string
}

type TVShow struct {
	Title    string
	Filename string
}

func (t TVShow) TableName() string {
	return "tvshows"
}
func (t TVShow) MediaFilename() string {
	return t.Filename
}
func (t TVShow) MediaTitle() string {
	return t.Filename
}

func (t TVShow) MediaType() string {
	return "tvshow"
}

type Movie struct {
	Title    string
	Filename string
}

func (m Movie) MediaFilename() string {
	return m.Filename
}
func (m Movie) MediaTitle() string {
	return m.Filename
}

func (m Movie) MediaType() string {
	return "movie"
}

type Other struct {
	Title    string
	Filename string
}

func (o Other) MediaFilename() string {
	return o.Filename
}
func (o Other) MediaTitle() string {
	return o.Filename
}

func (o Other) MediaType() string {
	return "other"
}

func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", dir+string(filepath.Separator)+"btpodcast.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	db.LogMode(true)
	db.AutoMigrate(&TVShow{}, &Movie{}, &Other{})
	return db
}

var db = initDB()
