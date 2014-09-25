package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
)

type TVShow struct {
	Name string
}

func (t TVShow) TableName() string {
	return "tvshows"
}

type Movie struct {
	Name string
}

type Other struct {
	Name string
}

func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", dir+string(filepath.Separator)+"btpodcast.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	db.SingularTable(true)
	db.AutoMigrate(&TVShow{}, &Movie{}, &Other{})
	return db
}

var db = initDB()
