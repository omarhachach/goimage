package main

import (
	"flag"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
)

var (
	configFilePath = flag.String("config", "./config.json", "-config <config-path>")
	debug          = flag.Bool("debug", false, "-debug <true/false>")
	config         *Config
)

func init() {
	flag.Parse()

	if *debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	var err error
	config, err = ParseConfig(*configFilePath)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse config.")
		return
	}

	err = os.MkdirAll(config.Directories.Template, 644)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create template directory.")
		return
	}

	err = os.MkdirAll(config.Directories.Public, 644)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create public directory.")
		return
	}

	err = os.MkdirAll(config.Directories.Image, 644)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create image directory.")
		return
	}
}

func main() {
	router := chi.NewRouter()

	router.Use(middleware.Timeout(30 * time.Second))
	router.Use(middleware.Recoverer)
	router.Use(middleware.DefaultCompress)
	router.Use(MaxBodySizeMiddleware)

	if *debug {
		router.Use(middleware.Logger)
	}

	templates, err := template.ParseGlob(config.Directories.Template + "*.html")
	if err != nil {
		logrus.WithError(err).Fatal("Failed to parse templates.")
		return
	}

	indexTemplate := templates.Lookup("index.html")
	if indexTemplate == nil {
		logrus.Fatal("Failed to parse index template (index.html).")
		return
	}

	viewTemplate := templates.Lookup("view.html")
	if viewTemplate == nil {
		logrus.Fatal("Failed to parse view template (view.html).")
		return
	}

	notFoundTemplate := templates.Lookup("404.html")
	if notFoundTemplate == nil {
		logrus.Fatal("Failed to parse 404 template (404.html).")
		return
	}

	router.Get("/", IndexHandler(indexTemplate))
	router.Get("/{id}/", ViewHandler(viewTemplate, notFoundTemplate))
	router.Post("/upload/", UploadHandler)
	FileServer(router, "/", http.Dir(config.Directories.Public))

	if config.CSRF.Enabled {
		http.ListenAndServe(":"+strconv.Itoa(config.Port), csrf.Protect(
			[]byte(config.CSRF.AuthKey),
			csrf.CookieName("_csrf"),
			csrf.FieldName("_csrf"),
			csrf.HttpOnly(config.CSRF.HTTPOnly),
			csrf.Secure(config.CSRF.Secure),
		)(router))
	} else {
		http.ListenAndServe(":"+strconv.Itoa(config.Port), router)
	}
}
