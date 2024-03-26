package main

import (
	"log"
	"os"

	"github.com/BurntSushi/toml"
)

var config struct {
	Address  string
	Database string
}

func init() {
	path := os.Getenv("NECHEGO_CONFIG")
	if path == "" {
		path = "config.toml"
	}

	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Fatal(err)
	}

	if config.Address == "" {
		config.Address = ":8080"
	}
}
