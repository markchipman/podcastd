package main

import (
	"fmt"
	"net/http"
)

func main() {
	initDirectory()
	updateDB()
	watchDownloads()
	http.HandleFunc("/", home)
	http.HandleFunc("/rss/movies", MoviesRSS)
	fs := http.FileServer(http.Dir(dir))
	http.Handle("/files/", http.StripPrefix("/files/", fs))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
