package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
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
	http.Handle("/feed/audio", auth(http.HandlerFunc(AudioFeed)))
	http.Handle("/feed/video", auth(http.HandlerFunc(VideoFeed)))

	movieFileServer := http.FileServer(http.Dir(config.Movies))
	http.Handle("/media/movies/", auth(http.StripPrefix("/media/movies/", movieFileServer)))

	tvshowFileServer := http.FileServer(http.Dir(config.TVShows))
	http.Handle("/media/tvshows/", auth(http.StripPrefix("/media/tvshows/", tvshowFileServer)))

	audioFileServer := http.FileServer(http.Dir(config.Audio))
	http.Handle("/media/audio/", auth(http.StripPrefix("/media/audio/", audioFileServer)))

	videoFileServer := http.FileServer(http.Dir(config.Video))
	http.Handle("/media/video/", auth(http.StripPrefix("/media/video/", videoFileServer)))

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}

//HandBrakeCLI -i test.avi -o output.m4v --preset="AppleTV 2"
