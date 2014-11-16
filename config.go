package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"strings"
)

type Config struct {
	Username string
	Password string
	Port     int
	Database string
	Movies   []string
	TVShows  []string
	Audio    []string
	Video    []string
}

func loadConfigFile() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// Convert tildes (~) in config paths to user home directory
	usr, _ := user.Current()
	config.Database = strings.Replace(config.Database, "~", usr.HomeDir, 1)
	for i, _ := range config.Movies {
		config.Movies[i] = strings.Replace(config.Movies[i], "~", usr.HomeDir, 1)
	}
	for i, _ := range config.TVShows {
		config.TVShows[i] = strings.Replace(config.TVShows[i], "~", usr.HomeDir, 1)
	}
	for i, _ := range config.Audio {
		config.Audio[i] = strings.Replace(config.Audio[i], "~", usr.HomeDir, 1)
	}
	for i, _ := range config.Video {
		config.Video[i] = strings.Replace(config.Video[i], "~", usr.HomeDir, 1)
	}

	return config
}

var config = loadConfigFile()
