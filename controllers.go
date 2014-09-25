package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "<h1>btPodcast</h1>")
	fmt.Fprint(w, "<h4>TV Shows</h4>")
	var tvshows []TVShow
	db.Find(&tvshows)
	for _, show := range tvshows {
		fmt.Fprint(w, show.Name)
		fmt.Fprint(w, "<br>")
	}
	fmt.Fprint(w, "<h4>Movies</h4>")
	var movies []Movie
	db.Find(&movies)
	for _, movie := range movies {
		fmt.Fprint(w, movie.Name)
		fmt.Fprint(w, "<br>")
	}
	fmt.Fprint(w, "<h4>Other</h4>")
	var other []Other
	db.Find(&other)
	for _, file := range other {
		fmt.Fprint(w, file.Name)
		fmt.Fprint(w, "<br>")
	}
}
