package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
)

func main() {
	updateDB()
	watchDir(config.Movies, ProcessMovie)
	watchDirs(config.TVShows, ProcessTVShow)
	watchDir(config.Audio, ProcessAudio)
	watchDir(config.Video, ProcessVideo)
	auth := httpauth.SimpleBasicAuth(config.Username, config.Password)
	http.Handle("/", auth(http.HandlerFunc(home)))
	http.Handle("/feed/movies", auth(http.HandlerFunc(MoviesRSS)))
	movieFileServer := http.FileServer(http.Dir(config.Movies))
	http.Handle("/media/movies/", auth(http.StripPrefix("/media/movies/", movieFileServer)))
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

//HandBrakeCLI -i test.avi -o output.m4v --preset="AppleTV 2"
