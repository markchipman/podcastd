package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
)

func main() {
	updateDB()
	watchDownloads()
	auth := httpauth.SimpleBasicAuth(config.Username, config.Password)
	http.Handle("/", auth(http.HandlerFunc(home)))
	http.Handle("/rss/movies", auth(http.HandlerFunc(MoviesRSS)))
	for i, _ := range config.Movies {
		prefix := fmt.Sprintf("/media/movies/%d/", i)
		http.Handle(prefix, auth(http.StripPrefix(prefix, http.FileServer(http.Dir(config.Movies[i])))))
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

//HandBrakeCLI -i test.avi -o output.m4v --preset="AppleTV 2"
