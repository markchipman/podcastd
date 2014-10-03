package main

import (
	"html/template"
	"net/http"
	"path/filepath"
)

var tmplDir = "templates" + string(filepath.Separator)
var templates = template.Must(template.ParseFiles(tmplDir + "home.html"))

func home(w http.ResponseWriter, r *http.Request) {
	var movies []Movie
	db.Find(&movies)
	var shows []string
	rows, _ := db.Raw("SELECT DISTINCT show_title FROM tvshows").Rows()
	defer rows.Close()
	for rows.Next() {
		var title string
		rows.Scan(&title)
		shows = append(shows, title)
	}
	var tvshows map[string]interface{}
	tvshows = make(map[string]interface{})
	for _, show := range shows {
		var episodes []TVShow
		db.Where("show_title = ?", show).Find(&episodes)
		tvshows[show] = episodes
	}
	data := map[string]interface{}{
		"movies":  movies,
		"tvshows": tvshows,
	}
	err := templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
