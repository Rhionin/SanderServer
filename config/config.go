package config

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

// Config is the config.
type Config struct {
	ProgressCheckInterval string `yaml:"progressCheckInterval"`
	Port                  string `yaml:"port"`
	ProgressTopic         string `yaml:"progresstopic"`
	SlackWebhookURL       string `yaml:"slackWebhookURL"`
	GithubUsername        string `yaml:"githubUsername"`
	GithubApiKey          string `yaml:"githubApiKey"`
}

// GetConfig returns the config object. Gofigure
func GetConfig(filePath string) Config {

	filename, _ := filepath.Abs(filePath)
	yamlFile, err := os.ReadFile(filename)

	if err != nil {
		panic(err)
	}

	var config Config

	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		panic(err)
	}

	if username := os.Getenv("GIT_USERNAME"); username != "" {
		config.GithubUsername = username
	}
	if apiKey := os.Getenv("GIT_API_KEY"); apiKey != "" {
		config.GithubApiKey = apiKey
	}

	return config
}
