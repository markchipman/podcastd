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
	Movies   string
	TVShows  string
	Audio    string
	Video    string
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
	config.Movies = strings.Replace(config.Movies, "~", usr.HomeDir, 1)
	config.TVShows = strings.Replace(config.TVShows, "~", usr.HomeDir, 1)
	config.Audio = strings.Replace(config.Audio, "~", usr.HomeDir, 1)
	config.Video = strings.Replace(config.Video, "~", usr.HomeDir, 1)

	return config
}

var config = loadConfigFile()
