package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func main() {

	// Save media to database
	d, err := os.Open(dir + string(filepath.Separator) + "tvshows")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		show := TVShow{}
		db.FirstOrCreate(&show, TVShow{Filename: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}
	d, err = os.Open(dir + string(filepath.Separator) + "movies")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err = d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		movie := Movie{}
		db.FirstOrCreate(&movie, Movie{Filename: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}
	d, err = os.Open(dir + string(filepath.Separator) + "other")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err = d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		other := Other{}
		db.FirstOrCreate(&other, Other{Filename: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}

	// Setup and start web server
	http.HandleFunc("/", home)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
