package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/omar-h/goimage/api"
	"github.com/omar-h/goimage/api/config"
	"github.com/sirupsen/logrus"
)

var cfgPath = flag.String("config", "./config.json", "-config <path>")

func init() {
	flag.Parse()
}

func main() {
	config, err := config.ReadConfig(*cfgPath)
	if err != nil {
		logrus.WithError(err).Warn("Failed to read config, using default.")
	}

	err = os.MkdirAll(config.FileUploadLocation, 0644)
	if err != nil && err != os.ErrExist {
		logrus.WithError(err).Fatal("Failed to create upload location.")
	}

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
	)

	r.Post("/upload", api.UploadHandler)

	logrus.Infof("Now starting the server on port %v.", config.Port)
	http.ListenAndServe(":"+strconv.Itoa(config.Port), r)
}
