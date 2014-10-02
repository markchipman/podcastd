package main

import (
	"fmt"
	"net/http"
)

func main() {
	initDirectory()
	updateDB()
	http.HandleFunc("/", home)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
