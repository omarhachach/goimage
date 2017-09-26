package main

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds the configuration options for the image server.
type Config struct {
	Port              int      `json:"port"`
	AllowedMIMETypes  []string `json:"allowed_mime_types"`
	AllowedExtensions []string `json:"allowed_extensions"`
	ImageNameLength   int      `json:"image_name_length"`
	MaxFileSize       int64    `json:"max_file_size"`
	CSRF              struct {
		Enabled  bool   `json:"enabled"`
		AuthKey  string `json:"32_byte_auth_key"`
		Secure   bool   `json:"secure"`
		HTTPOnly bool   `json:"httpOnly"`
	} `json:"csrf"`
	Directories struct {
		Image    string `json:"image"`
		Template string `json:"template"`
		Public   string `json:"public"`
	} `json:"directories"`
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
