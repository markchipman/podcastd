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
	tmplDir+"tvseries.xml",
	tmplDir+"audio.xml",
	tmplDir+"video.xml",
))

func Home(w http.ResponseWriter, r *http.Request) {
	var movies []Media
	db.Where(Media{Type: "movie"}).Order("title").Find(&movies)
	type Series struct {
		Title string
		Slug  string
	}
	var tvshows map[string][]Media
	tvshows = make(map[string][]Media)
	rows, _ := db.Raw("SELECT DISTINCT title FROM media WHERE type = ? ORDER BY title", "tvshow").Rows()
	defer rows.Close()
	for rows.Next() {
		var title string
		rows.Scan(&title)
		slug := strings.ToLower(strings.Replace(title, " ", "-", -1))
		var shows []Media
		db.Where(Media{Type: "tvshow", Title: title}).Order("season, episode").Find(&shows)
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
	row := db.Raw("SELECT created_at FROM media WHERE type = ? ORDER BY created_at DESC LIMIT 1;", "movie").Row()
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
	var lastUpdate time.Time
	row := db.Raw("SELECT created_at FROM media WHERE type = ? ORDER BY created_at DESC LIMIT 1;", "tvshow").Row()
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

func TVSeriesFeed(w http.ResponseWriter, r *http.Request) {
	_, slug := path.Split(r.URL.Path)
	title := strings.Title(strings.Replace(slug, "-", " ", -1))
	var episodes []Media
	db.Where(Media{Type: "tvshow", Title: title}).Find(&episodes)
	var lastUpdate time.Time
	row := db.Raw("SELECT created_at FROM media WHERE type = ? AND title = ? ORDER BY created_at DESC LIMIT 1;", "tvshow", title).Row()
	row.Scan(&lastUpdate)
	data := map[string]interface{}{
		"lastUpdate": lastUpdate.Format(time.RFC1123),
		"host":       r.Host,
		"title":      episodes[0].Title,
		"desc":       episodes[0].Desc,
		"poster":     episodes[0].Poster,
		"episodes":   episodes,
	}
	err := xml.ExecuteTemplate(w, "tvseries.xml", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AudioFeed(w http.ResponseWriter, r *http.Request) {
	var audio []Media
	db.Where(Media{Type: "audio"}).Find(&audio)
	row := db.Raw("SELECT created_at FROM media WHERE type = ? ORDER BY created_at DESC LIMIT 1;", "audio").Row()
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
	row := db.Raw("SELECT created_at FROM media WHERE type = ? ORDER BY created_at DESC LIMIT 1;", "video").Row()
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
