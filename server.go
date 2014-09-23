package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Config struct {
	Port int
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "btPodcast")
}

func main() {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
	}
	http.HandleFunc("/", handler)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), nil)
}
