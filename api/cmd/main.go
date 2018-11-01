package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/omar-h/goimage/api"
)

func main() {
	r := chi.NewRouter()

	r.Use(
		middleware.Recoverer,
	)

	r.Post("/upload", api.UploadHandler)

	http.ListenAndServe(":8080", r)
}
