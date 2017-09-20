package main

import "io/ioutil"
import "encoding/json"

// Config holds the configuration options for the image server.
type Config struct {
	Port              int      `json:"port"`
	Secure            bool     `json:"secure"`
	AuthKey           string   `json:"32-byte-auth-key"`
	AllowedMimeTypes  []string `json:"allowed-mime-types"`
	AllowedExtensions []string `json:"allowed-extensions"`
	ImageNameLength   int      `json:"image-name-length"`
	MaxFileSize       int64    `json:"max-file-size"`
	ImageDirectory    string   `json:"image-directory"`
	TemplateDirectory string   `json:"template-directory"`
	PublicDirectory   string   `json:"public-directory"`
	ImageURL          string   `json:"image-url"`
	CSRF              bool     `json:"csrf"`
}

// ParseConfig parses the passed filepath and returns a new Config.
func ParseConfig(filepath string) (config *Config, err error) {
	configFile, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	parsedConfig := &Config{}
	err = json.Unmarshal(configFile, parsedConfig)
	if err != nil {
		return nil, err
	}

	return parsedConfig, nil
}
