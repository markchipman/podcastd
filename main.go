package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
)

func main() {
	go func() {
		updateDB()
		for _, dir := range config.Media {
			watchDir(dir)
		}
	}()

	auth := httpauth.SimpleBasicAuth(config.Username, config.Password)
	http.Handle("/", auth(http.HandlerFunc(Home)))
	http.Handle("/media/", http.HandlerFunc(MediaFile))
	http.Handle("/feed/movies", http.HandlerFunc(MovieFeed))
	http.Handle("/feed/tvshows", http.HandlerFunc(TVShowFeed))
	http.Handle("/feed/tvshows/", http.HandlerFunc(TVSeriesFeed))

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
