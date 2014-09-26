package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var tmplDir = "templates" + string(filepath.Separator)
var templates = template.Must(template.ParseFiles(tmplDir + "home.html"))

func home(w http.ResponseWriter, r *http.Request) {
	data := make(map[string][]Media)
	var tvshows []TVShow
	db.Find(&tvshows)
	data["TV Shows"] = make([]Media, len(tvshows))
	for i, v := range tvshows {
		data["TV Shows"][i] = Media(v)
	}
	var movies []Movie
	db.Find(&movies)
	data["Movies"] = make([]Media, len(movies))
	for i, v := range movies {
		data["Movies"][i] = Media(v)
	}
	var other []Other
	db.Find(&other)
	data["Other Media"] = make([]Media, len(other))
	for i, v := range other {
		data["Other Media"][i] = Media(v)
	}
	err := templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
