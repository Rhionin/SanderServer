package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path/filepath"
)

// Config is the config.
type Config struct {
	GoogleAPIKey          string `yaml:"googleAPIKey"`
	ProgressCheckInterval string `yaml:"progressCheckInterval"`
}

// GetConfig returns the config object. Gofigure
func GetConfig() Config {

	filename, _ := filepath.Abs("./config/config.yaml")
	yamlFile, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	return config
}
