package main

import (
	"html/template"
	"net/http"
	"path/filepath"
	xtemplate "text/template"
	"time"
)

var tmplDir = "templates" + string(filepath.Separator)
var templates = template.Must(template.ParseFiles(tmplDir + "home.html"))
var xml = xtemplate.Must(xtemplate.ParseFiles(
	tmplDir+"movies.xml",
	tmplDir+"tvshows.xml",
	tmplDir+"audio.xml",
	tmplDir+"video.xml",
))

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
	var audio []Audio
	db.Find(&audio)
	var video []Video
	db.Find(&video)
	data := map[string]interface{}{
		"movies":  movies,
		"tvshows": tvshows,
		"audio":   audio,
		"video":   video,
		"host":    r.Host,
	}
	err := templates.ExecuteTemplate(w, "home.html", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func MovieFeed(w http.ResponseWriter, r *http.Request) {
	var movies []Movie
	db.Find(&movies)
	row := db.Raw("SELECT added FROM movies ORDER BY added DESC LIMIT 1;").Row()
	var lastUpdate time.Time
	row.Scan(&lastUpdate)
	data := map[string]interface{}{
		"lastUpdate": lastUpdate.Format(time.RFC1123),
		"host":       r.Host,
		"movies":     movies,
	}
	err := xml.ExecuteTemplate(w, "movies.xml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func TVShowFeed(w http.ResponseWriter, r *http.Request) {
	var tvshows []TVShow
	db.Find(&tvshows)
	row := db.Raw("SELECT aired FROM tvshows ORDER BY aired DESC LIMIT 1;").Row()
	var lastUpdate time.Time
	row.Scan(&lastUpdate)
	data := map[string]interface{}{
		"lastUpdate": lastUpdate.Format(time.RFC1123),
		"host":       r.Host,
		"tvshows":    tvshows,
	}
	err := xml.ExecuteTemplate(w, "tvshows.xml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AudioFeed(w http.ResponseWriter, r *http.Request) {
	var audio []Audio
	db.Find(&audio)
	row := db.Raw("SELECT added FROM audio ORDER BY added DESC LIMIT 1;").Row()
	var lastUpdate time.Time
	row.Scan(&lastUpdate)
	data := map[string]interface{}{
		"lastUpdate": lastUpdate.Format(time.RFC1123),
		"host":       r.Host,
		"audio":      audio,
	}
	err := xml.ExecuteTemplate(w, "audio.xml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func VideoFeed(w http.ResponseWriter, r *http.Request) {
	var video []Video
	db.Find(&video)
	row := db.Raw("SELECT added FROM video ORDER BY added DESC LIMIT 1;").Row()
	var lastUpdate time.Time
	row.Scan(&lastUpdate)
	data := map[string]interface{}{
		"lastUpdate": lastUpdate.Format(time.RFC1123),
		"host":       r.Host,
		"video":      video,
	}
	err := xml.ExecuteTemplate(w, "video.xml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
