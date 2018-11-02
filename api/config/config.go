package config

import (
	"encoding/json"
	"io/ioutil"
)

// Config holds the server param configurations.
type Config struct {
	Port               int      `json:"port,omitempty"`
	MaxFileSize        int      `json:"max_file_size,omitempty"`
	NameLength         int      `json:"name_length,omitempty"`
	FileBufferSize     int64    `json:"file_buffer_size,omitempty"`
	FileUploadLocation string   `json:"file_upload_location,omitempty"`
	AllowedExtensions  []string `json:"allowed_extensions,omitempty"`
	AllowedMIMETypes   []string `json:"allowed_mime_types,omitempty"`
}

// DefaultConfig is the config used if one hasn't been specified.
var DefaultConfig = &Config{
	Port:               8080,
	MaxFileSize:        20 << 18, // 5 MB
	NameLength:         6,
	FileBufferSize:     20 << 16, // 1.25 MB
	FileUploadLocation: "./img/",
	AllowedExtensions: []string{
		".png",
		".jpg",
		".jpeg",
		".jiff",
		".ico",
		".gif",
		".tif",
		".webp",
	},
	AllowedMIMETypes: []string{
		"image/x-icon",
		"image/jpeg",
		"image/pjpeg",
		"image/png",
		"image/tiff",
		"image/x-tiff",
		"image/webp",
		"image/gif",
	},
}

// Cfg is the globally available config.
var Cfg = DefaultConfig

// ReadConfig will read a config from a given path.
func ReadConfig(path string) (*Config, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return Cfg, err
	}

	cfg := &Config{}
	err = json.Unmarshal(f, cfg)
	if err != nil {
		return Cfg, err
	}

	Cfg = cfg

	return cfg, nil
}
