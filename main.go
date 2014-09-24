package main

import (
	"encoding/json"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// Load configuration
type Config struct {
	Home string
	Port int
}

func LoadConfig() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return config
}

var config = LoadConfig()

// Define models
type TVShow struct {
	Name string
}

func (t TVShow) TableName() string {
	return "tvshows"
}

type Movie struct {
	Name string
}

type Other struct {
	Name string
}

// Setup btpodcast home directory
func btPodcastDir() string {
	home := config.Home + string(filepath.Separator)
	if strings.Contains(home, "~/") {
		usr, _ := user.Current()
		home = strings.Replace(home, "~/", usr.HomeDir+string(filepath.Separator), 1)
	}
	os.MkdirAll(home+string(filepath.Separator)+"tvshows", 0700)
	os.MkdirAll(home+string(filepath.Separator)+"movies", 0700)
	os.MkdirAll(home+string(filepath.Separator)+"other", 0700)
	return home
}

var home = btPodcastDir()

// Initialize database
func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", home+string(filepath.Separator)+"btpodcast.db")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	db.SingularTable(true)
	db.AutoMigrate(&TVShow{}, &Movie{}, &Other{})
	return db
}

var db = initDB()

func handler(w http.ResponseWriter, r *http.Request) {
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

func main() {

	// Save media to database
	d, err := os.Open(home + string(filepath.Separator) + "tvshows")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		show := TVShow{}
		db.FirstOrCreate(&show, TVShow{Name: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}
	d, err = os.Open(home + string(filepath.Separator) + "movies")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err = d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		movie := Movie{}
		db.FirstOrCreate(&movie, Movie{Name: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}
	d, err = os.Open(home + string(filepath.Separator) + "other")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer d.Close()
	files, err = d.Readdir(-1)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, file := range files {
		other := Other{}
		db.FirstOrCreate(&other, Other{Name: file.Name()})
		fmt.Println(file.Name(), file.Size(), "bytes")
	}

	// Setup and start web server
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
