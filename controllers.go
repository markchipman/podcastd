package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var tmplDir = "templates" + string(filepath.Separator)
var templates = template.Must(template.ParseFiles(tmplDir + "home.html"))

func home(w http.ResponseWriter, r *http.Request) {
	var tvshows []TVShow
	db.Find(&tvshows)
	var movies []Movie
	db.Find(&movies)
	var others []Other
	db.Find(&others)
	data := map[string]interface{}{
		"tvshows": tvshows,
		"movies":  movies,
		"others":  others,
	}
	err := templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
