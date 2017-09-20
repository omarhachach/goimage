package main

import (
	"flag"
	"html/template"
	"os"

	"github.com/buaazp/fasthttprouter"
	"github.com/sirupsen/logrus"
)

var (
	// configFilePath holds the
	configFilePath = flag.String("config", "./config.json", "The JSON config filepath.")
	config         *Config
)

func init() {
	flag.Parse()

	var err error
	config, err = ParseConfig(*configFilePath)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to parse config.")
		return
	}

	err = os.MkdirAll(config.TemplateDirectory, 644)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create template directory.")
		return
	}

	err = os.MkdirAll(config.PublicDirectory, 644)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create public directory.")
		return
	}

	err = os.MkdirAll(config.ImageDirectory, 644)
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to create image directory.")
		return
	}
}

func main() {
	router := fasthttprouter.New()

	templates, err := template.ParseGlob(config.TemplateDirectory + "*.html")
	if err != nil {
		logrus.WithField("error", err).Fatal("Failed to parse templates.")
		return
	}

	index := templates.Lookup("index.html")
	router.GET("/", indexHandler(index))
}
