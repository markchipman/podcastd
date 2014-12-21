package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/garfunkel/go-tvdb"
	_ "github.com/mattn/go-sqlite3"
	"github.com/ryanss/gorm"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var ValidFileType = map[string]bool{
	".m4a": true,
	".m4v": true,
	".mp3": true,
	".mp4": true,
}

type Media struct {
	Id           int
	Type         string
	Path         string
	Filename     string
	Size         int64
	Title        string
	Desc         string
	Runtime      int
	Genres       string
	Poster       string
	Season       int
	Episode      int
	EpisodeTitle string
	EpisodeDesc  string
	Released     time.Time
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}

func (m Media) TableName() string {
	return "media"
}

func (m Media) PubDate() string {
	if m.Released.Year() > 1900 {
		return m.Released.Format(time.RFC1123)
	}
	return m.CreatedAt.Format(time.RFC1123)
}

func (m Media) Aired() string {
	return m.Released.Format("01/02/2006")
}

func (m Media) MediaURL(host string) string {
	Url, _ := url.Parse(fmt.Sprintf("http://%s/media/%d/%s", host, m.Id, m.Filename))
	return Url.String()
}

func (m Media) TitleSlug() string {
	return strings.ToLower(strings.Replace(m.Title, " ", "-", -1))
}

func (m Media) S00E00() string {
	return fmt.Sprintf("S%02dE%02d", m.Season, m.Episode)
}

func ProcessFile(fp string, timestamp time.Time) {
	file, _ := os.Stat(fp)
	fp, _ = path.Split(fp)
	media := Media{
		Path:     fp,
		Filename: file.Name(),
	}
	db.Where(media).FirstOrCreate(&media)
	media.Size = file.Size()

	if media.Type != "" {
		db.Save(&media)
		return
	}

	// Audio
	if filepath.Ext(file.Name()) == ".mp3" {
		media.Type = "audio"
	}

	// TV Show
	re := regexp.MustCompile("S[0-9]{2}E[0-9]{2}")
	info := re.FindString(media.Filename)
	if info != "" {
		media.Season, _ = strconv.Atoi(info[1:3])
		media.Episode, _ = strconv.Atoi(info[4:])
		i := re.FindStringIndex(media.Filename)
		media.Title = strings.Replace(media.Filename[0:i[0]-1], ".", " ", -1)
		media.Type = "tvshow"
		media.ScrapeTVShow()
	}

	// Movie
	if media.Type == "" {
		filename := []byte(media.Filename)
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
		media.Title = strings.Replace(string(filename[0:index]), ".", " ", -1)
		if iYear != nil {
			media.Released, _ = time.Parse("2006", string(filename[iYear[0]+1:iYear[1]])+"-01-01")
		}
		media.ScrapeMovie()
		if media.Desc != "" {
			media.Type = "movie"
		}
	}

	// Video
	if media.Type == "" {
		media.Type = "video"
	}

	db.Save(&media)
}

func (m *Media) ScrapeMovie() {
	searchURL := "https://www.themoviedb.org/search/movie?query=" + m.Title
	searchURL = strings.Replace(searchURL, " ", "%20", -1)
	doc, _ := goquery.NewDocument(searchURL)
	s := doc.Find("ul.movie li").First()
	s = s.Find("a").First()
	link, _ := s.Attr("href")
	doc, _ = goquery.NewDocument("https://www.themoviedb.org" + link)
	s = doc.Find("#overview").First()
	m.Desc = s.Text()
	doc.Find("#genres span").Each(func(i int, s *goquery.Selection) {
		m.Genres = m.Genres + s.Text() + ", "
	})
	if m.Genres != "" {
		m.Genres = m.Genres[:len(m.Genres)-2]
	}
	s = doc.Find("span[itemprop=datePublished]").First()
	if s.Text() == "" {
		s = doc.Find("#release_date_list span").First()
	}
	m.Released, _ = time.ParseInLocation("2006-01-02", s.Text(), time.Local)
	s = doc.Find("a.poster").First()
	m.Poster, _ = s.Find("img").Attr("src")
	runtime, _ := strconv.ParseInt(doc.Find("#runtime").Text(), 10, 0)
	m.Runtime = int(runtime)
}

func (m *Media) ScrapeTVShow() {
	const APIKey = "E01489F781B562D8"
	seriesList, _ := tvdb.SearchSeries(m.Title, 1)
	series := seriesList.Series[0]
	series.GetDetail()
	m.Desc = series.Overview
	runtime, _ := strconv.ParseInt(series.Runtime, 10, 0)
	m.Runtime = int(runtime)
	m.Genres = strings.Join(series.Genre, ", ")
	m.Poster = "http://thetvdb.com/banners/_cache/" + series.Poster
	episode := series.Seasons[uint64(m.Season)][m.Episode-1]
	m.EpisodeTitle = episode.EpisodeName
	m.EpisodeDesc = episode.Overview
	m.Released, _ = time.ParseInLocation("2006-01-02", episode.FirstAired, time.Local)
}

func initDB() gorm.DB {
	db, err := gorm.Open("sqlite3", config.Database)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	db.DB()
	//db.LogMode(true)
	db.AutoMigrate(&Media{})
	return db
}

var db = initDB()

func updateDB() {
	timestamp := time.Now().Local()

	for _, dir := range config.Media {
		filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
			if ValidFileType[filepath.Ext(f.Name())] {
				ProcessFile(path, timestamp)
			}
			return nil
		})
	}

	// Soft delete records that were not found
	db.Where("updated_at < ?", timestamp).Delete(Media{})
}
