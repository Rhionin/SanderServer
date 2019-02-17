package config

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Config is the config.
type Config struct {
	ProgressCheckInterval string `yaml:"progressCheckInterval"`
	Port                  string `yaml:"port"`
	ProgressTopic         string `yaml:"progresstopic"`
}

// GetConfig returns the config object. Gofigure
func GetConfig(filePath string) Config {

	filename, _ := filepath.Abs(filePath)
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
