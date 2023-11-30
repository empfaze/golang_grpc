package utils

import (
	"flag"
	"os"
)

func GetConfigPath() string {
	var path string

	flag.StringVar(&path, "config", "", "Path to config file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}
