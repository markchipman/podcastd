package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type Config struct {
	Home string
	Port int
}

func loadConfigFile() Config {
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Config{}
	err := decoder.Decode(&config)
	if err != nil {
		fmt.Println("Error:", err)
	}
	return config
}

var config = loadConfigFile()

func initDirectory() string {
	home := config.Home + string(filepath.Separator)
	if strings.Contains(home, "~/") {
		usr, _ := user.Current()
		home = strings.Replace(home, "~/", usr.HomeDir+string(filepath.Separator), 1)
	}
	os.MkdirAll(home+string(filepath.Separator)+"tvshows", 0700)
	os.MkdirAll(home+string(filepath.Separator)+"movies", 0700)
	os.MkdirAll(home+string(filepath.Separator)+"other", 0700)
	return home
}

var dir = initDirectory()
