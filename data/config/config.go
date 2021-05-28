package config

import (
	"carousell-gobot/models"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Config - contains all the configuration of the program
var Config = models.DefaultConfig()

// Load - load config
func Load(path string) error {
	path, err := getPath(path)
	if err != nil {
		return err
	}

	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(raw, Config)
	if err != nil {
		return err
	}

	return nil
}

func getPath(path string) (string, error) {
	if strings.HasPrefix(path, "/") {
		return path, nil
	}
	absPath, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	absPath += "/" + path
	return absPath, nil
}
