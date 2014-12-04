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
	Media    []string
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
	for i, dir := range config.Media {
		config.Media[i] = strings.Replace(dir, "~", usr.HomeDir, 1)
	}
	return config
}

var config = loadConfigFile()
