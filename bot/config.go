package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var config struct {
	Token    string
	Database string
	Assets   string
}

func init() {
	path := os.Getenv("NECHEGO_CONFIG")
	if path == "" {
		path = "config.toml"
	}

	_, err := toml.DecodeFile(path, &config)
	if err != nil {
		log.Fatal(err)
	}

	if config.Token == "" {
		log.Fatal("config: token must be set")
	}
	if config.Database == "" {
		log.Fatal("config: database must be set")
	}
	if config.Assets == "" {
		config.Assets = "assets"
	}
}
