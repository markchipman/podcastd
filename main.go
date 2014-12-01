package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
	"strings"
)

func main() {
	updateDB()
	watchDir(config.Movies, ProcessMovie)
	watchDirs(config.TVShows, ProcessTVShow)
	watchDir(config.Audio, ProcessAudio)
	watchDir(config.Video, ProcessVideo)

	auth := httpauth.SimpleBasicAuth(config.Username, config.Password)
	http.Handle("/", auth(http.HandlerFunc(home)))
	http.Handle("/feed/movies", auth(http.HandlerFunc(MovieFeed)))
	http.Handle("/feed/tvshows", auth(http.HandlerFunc(TVShowFeed)))
	http.Handle("/feed/audio", auth(http.HandlerFunc(AudioFeed)))
	http.Handle("/feed/video", auth(http.HandlerFunc(VideoFeed)))

	http.Handle("/movies/", http.HandlerFunc(MovieFile))

	tvshowFileServer := http.FileServer(http.Dir(config.TVShows))
	http.Handle("/media/tvshows/", auth(http.StripPrefix("/media/tvshows/", tvshowFileServer)))
	rows, _ := db.Raw("SELECT DISTINCT show_title FROM tvshows").Rows()
	defer rows.Close()
	for rows.Next() {
		var title string
		rows.Scan(&title)
		slug := strings.ToLower(strings.Replace(title, " ", "", -1))
		http.Handle("/media/tvshows/"+slug+"/", http.StripPrefix("/media/tvshows/"+slug+"/", http.FileServer(http.Dir(config.TVShows+"/"+title+"/"))))
	}

	audioFileServer := http.FileServer(http.Dir(config.Audio))
	http.Handle("/media/audio/", http.StripPrefix("/media/audio/", audioFileServer))

	videoFileServer := http.FileServer(http.Dir(config.Video))
	http.Handle("/media/video/", http.StripPrefix("/media/video/", videoFileServer))

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

//HandBrakeCLI -i test.avi -o output.m4v --preset="AppleTV 2"
