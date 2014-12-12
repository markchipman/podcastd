package main

import (
	"html/template"
	"net/http"
	"path"
	"path/filepath"
	"strconv"
	"strings"
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

func Home(w http.ResponseWriter, r *http.Request) {
	var movies []Media
	db.Where(Media{Type: "movie"}).Find(&movies)
	type Series struct {
		Title string
		Slug  string
	}
	var tvshows map[string][]Media
	tvshows = make(map[string][]Media)
	rows, _ := db.Raw("SELECT DISTINCT title FROM media WHERE type = ?", "tvshow").Rows()
	defer rows.Close()
	for rows.Next() {
		var title string
		rows.Scan(&title)
		slug := strings.ToLower(strings.Replace(title, " ", "-", -1))
		var shows []Media
		db.Where(Media{Type: "tvshow", Title: title}).Find(&shows)
		tvshows[slug] = shows
	}
	var audio []Media
	db.Where(Media{Type: "audio"}).Find(&audio)
	var video []Media
	db.Where(Media{Type: "video"}).Find(&video)
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

func MediaFile(w http.ResponseWriter, r *http.Request) {
	p, filename := path.Split(r.URL.Path)
	_, mediaId := path.Split(p[:len(p)-1])
	id, _ := strconv.ParseInt(mediaId, 10, 0)
	media := Media{Id: int(id), Filename: filename}
	db.Where(&media).First(&media)
	http.ServeFile(w, r, media.Path+string(filepath.Separator)+media.Filename)
}

func MovieFeed(w http.ResponseWriter, r *http.Request) {
	var movies []Media
	db.Where(Media{Type: "movie"}).Find(&movies)
	row := db.Raw("SELECT created_at FROM media ORDER BY created_at DESC LIMIT 1;").Row()
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
	var tvshows []Media
	db.Where(Media{Type: "tvshow"}).Find(&tvshows)
	row := db.Raw("SELECT created_at FROM media ORDER BY created_at DESC LIMIT 1;").Row()
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
	var audio []Media
	db.Find(&audio)
	row := db.Raw("SELECT created_at FROM media ORDER BY created_at DESC LIMIT 1;").Row()
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
	var video []Media
	db.Find(&video)
	row := db.Raw("SELECT created_at FROM media ORDER BY created_at DESC LIMIT 1;").Row()
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
