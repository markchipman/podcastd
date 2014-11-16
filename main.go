package main

import (
	"fmt"
	"net/http"
)

func main() {
	updateDB()
	watchDownloads()
	http.HandleFunc("/", home)
	http.HandleFunc("/rss/movies", MoviesRSS)
	for i, _ := range config.Movies {
		prefix := fmt.Sprintf("/media/movies/%d/", i)
		http.Handle(prefix, http.StripPrefix(prefix, http.FileServer(http.Dir(config.Movies[i]))))
	}
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

//HandBrakeCLI -i test.avi -o output.m4v --preset="AppleTV 2"
