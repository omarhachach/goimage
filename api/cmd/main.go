package main

import (
	"flag"
	"net/http"
	"os"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/omarhachach/goimage/api"
	"github.com/omarhachach/goimage/api/config"
	"github.com/sirupsen/logrus"
)

var cfgPath = flag.String("config", "./config.json", "-config <path>")

func init() {
	flag.Parse()
}

func main() {
	configStr, err := config.ReadConfig(*cfgPath)
	if err != nil {
		logrus.WithError(err).Warn("Failed to read config, using default.")
	}

	err = os.MkdirAll(configStr.FileUploadLocation, 0644)
	if err != nil && err != os.ErrExist {
		logrus.WithError(err).Fatal("Failed to create upload location: " + configStr.FileUploadLocation)
	}

	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
	)

	r.Post("/upload", api.UploadHandler)

	logrus.Infof("Now starting the server on port %v.", configStr.Port)
	http.ListenAndServe(":"+strconv.Itoa(configStr.Port), r)
}
