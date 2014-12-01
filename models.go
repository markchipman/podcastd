package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ryanss/gorm"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var validFileType = map[string]bool{
	".m4a": true,
	".m4v": true,
	".mp3": true,
	".mp4": true,
}

type Movie struct {
	Id        int
	Filename  string
	Size      int64
	Title     string
	Year      int
	Desc      string
	Genres    string
	Poster    string
	Runtime   int
	Added     time.Time
	Timestamp time.Time
}

func (m Movie) PubDate() string {
	return m.Added.Format(time.RFC1123)
}

func (m Movie) MediaURL(host string) string {
	Url, _ := url.Parse(fmt.Sprintf("http://%s/movies/%d", host, m.Id))
	return Url.String()
}

func (m *Movie) Parse() {
	filename := []byte(m.Filename)
	index := len(filename)
	reYear := regexp.MustCompile("\\.[0-9]{4}")
	iYear := reYear.FindIndex(filename)
	if iYear != nil && iYear[0] < index {
		index = iYear[0]
	}
	reExt := regexp.MustCompile("\\.[a-z0-9]+$")
	iExt := reExt.FindIndex(filename)
	if iExt != nil && iExt[0] < index {
		index = iExt[0]
	}
	m.Title = strings.Replace(string(filename[0:index]), ".", " ", -1)
	if iYear != nil {
		year, _ := strconv.ParseInt(string(filename[iYear[0]+1:iYear[1]]), 10, 0)
		m.Year = int(year)
	}
}

func (m *Movie) Scrape() {
	searchURL := "https://www.themoviedb.org/search?query=" + m.Title
	searchURL = strings.Replace(searchURL, " ", "%20", -1)
	doc, _ := goquery.NewDocument(searchURL)
	s := doc.Find("ul.movie li").First()
	s = s.Find("a").First()
	link, _ := s.Attr("href")
	fmt.Println(link)
	doc, _ = goquery.NewDocument("https://www.themoviedb.org" + link)
	s = doc.Find("#overview").First()
	m.Desc = s.Text()
	m.Genres = ""
	doc.Find("#genres span").Each(func(i int, s *goquery.Selection) {
		m.Genres = m.Genres + s.Text() + ", "
	})
	m.Genres = m.Genres[:len(m.Genres)-2]
	s = doc.Find("a.poster").First()
	m.Poster, _ = s.Find("img").Attr("src")
	runtime, _ := strconv.ParseInt(doc.Find("#runtime").Text(), 10, 0)
	m.Runtime = int(runtime)
}

func ProcessMovie(file os.FileInfo, timestamp time.Time) {
	if validFileType[filepath.Ext(file.Name())] {
		movie := Movie{
			Filename: file.Name(),
			Size:     file.Size(),
		}
		movie.Parse()
		db.Where(movie).Assign(Movie{Timestamp: timestamp}).FirstOrCreate(&movie)
		t, _ := time.Parse("2006-01-02", "1900-01-01")
		if t.After(movie.Added) {
			movie.Added = time.Now()
			movie.Scrape()
			db.Save(&movie)
		}
	}
}

type TVShow struct {
	Id           int
	Filename     string
	Size         int64
	ShowTitle    string
	Season       int
	Episode      int
	EpisodeTitle string
	EpisodeDesc  string
	Aired        time.Time
	Timestamp    time.Time
}

func (t TVShow) TableName() string {
	return "tvshows"
}

func (t TVShow) Slug() string {
	return strings.ToLower(strings.Replace(t.ShowTitle, " ", "", -1))
}

func (t TVShow) S00E00() string {
	return fmt.Sprintf("S%02dE%02d", t.Season, t.Episode)
}

func (t TVShow) MediaURL(host string) string {
	Url, _ := url.Parse(fmt.Sprintf("http://%s/media/tvshows/%s/%s", host, t.ShowTitle, t.Filename))
	return Url.String()
}

func (t *TVShow) Parse() {
	re := regexp.MustCompile("S[0-9]{2}E[0-9]{2}")
	info := re.FindString(t.Filename)
	t.Season, _ = strconv.Atoi(info[1:3])
	t.Episode, _ = strconv.Atoi(info[4:])
}

func ProcessTVShow(dir string, file os.FileInfo, timestamp time.Time) {
	if validFileType[filepath.Ext(file.Name())] {
		show := TVShow{
			ShowTitle: dir,
			Filename:  file.Name(),
			Size:      file.Size(),
		}
		show.Parse()
		db.Where(show).Assign(TVShow{Timestamp: timestamp}).FirstOrCreate(&show)
	}
}

type Audio struct {
	Id        int
	Filename  string
	Size      int64
	Added     time.Time
	Timestamp time.Time
}

func (a Audio) PubDate() string {
	return a.Added.Format(time.RFC1123)
}

func (a Audio) MediaURL(host string) string {
	Url, _ := url.Parse(fmt.Sprintf("http://%s/media/audio/%s", host, a.Filename))
	return Url.String()
}

func ProcessAudio(file os.FileInfo, timestamp time.Time) {
	if validFileType[filepath.Ext(file.Name())] {
		audio := Audio{}
		db.Where(Audio{
			Filename: file.Name(),
			Size:     file.Size(),
		}).Assign(Audio{Timestamp: timestamp}).FirstOrCreate(&audio)
	}
}

type Video struct {
	Id        int
	Filename  string
	Size      int64
	Added     time.Time
	Timestamp time.Time
}

func (v Video) PubDate() string {
	return v.Added.Format(time.RFC1123)
}

func (v Video) MediaURL(host string) string {
	Url, _ := url.Parse(fmt.Sprintf("http://%s/media/video/%s", host, v.Filename))
	return Url.String()
}

func ProcessVideo(file os.FileInfo, timestamp time.Time) {
	if validFileType[filepath.Ext(file.Name())] {
		video := Video{}
		db.Where(Video{
			Filename: file.Name(),
			Size:     file.Size(),
		}).Assign(Video{Timestamp: timestamp}).FirstOrCreate(&video)
	}
}

func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", config.Database)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	//db.LogMode(true)
	db.AutoMigrate(&Movie{}, &TVShow{}, &Audio{}, &Video{})
	return db
}

var db = initDB()

func updateDB() {
	timestamp := time.Now().Local()

	d, _ := os.Open(config.Movies)
	defer d.Close()
	files, _ := d.Readdir(-1)
	for _, file := range files {
		ProcessMovie(file, timestamp)
	}
	db.Exec("UPDATE movies SET added=datetime(?, 'localtime') WHERE added < '1990-01-01';", timestamp)

	d, _ = os.Open(config.TVShows)
	defer d.Close()
	files, _ = d.Readdir(-1)
	for _, file := range files {
		if file.IsDir() {
			e, _ := os.Open(config.TVShows + string(filepath.Separator) + file.Name())
			defer e.Close()
			episodes, _ := e.Readdir(-1)
			for _, episode := range episodes {
				ProcessTVShow(file.Name(), episode, timestamp)
			}
		}
	}

	d, _ = os.Open(config.Audio)
	defer d.Close()
	files, _ = d.Readdir(-1)
	for _, file := range files {
		ProcessAudio(file, timestamp)
	}
	db.Exec("UPDATE audios SET added=datetime(?, 'localtime') WHERE added < '1990-01-01';", timestamp)

	d, _ = os.Open(config.Video)
	defer d.Close()
	files, _ = d.Readdir(-1)
	for _, file := range files {
		ProcessVideo(file, timestamp)
	}
	db.Exec("UPDATE videos SET added=datetime(?, 'localtime') WHERE added < '1990-01-01';", timestamp)

	// Remove records from database that were not found
	db.Where("timestamp <> ?", timestamp).Delete(Movie{})
	db.Where("timestamp <> ?", timestamp).Delete(TVShow{})
	db.Where("timestamp <> ?", timestamp).Delete(Audio{})
	db.Where("timestamp <> ?", timestamp).Delete(Video{})
}
