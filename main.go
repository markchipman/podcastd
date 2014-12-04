package main

import (
	"fmt"
	"github.com/goji/httpauth"
	"net/http"
)

func main() {
	updateDB()
	for _, dir := range config.Media {
		watchDir(dir)
	}

	auth := httpauth.SimpleBasicAuth(config.Username, config.Password)
	http.Handle("/", auth(http.HandlerFunc(home)))
	http.Handle("/media/", http.HandlerFunc(MediaFile))

	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
