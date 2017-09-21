package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/omar-h/goimage"
	"github.com/sirupsen/logrus"
)

// FileServer serves the fileserver for the static files and images.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}

	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}

// IndexHandler serves the index at the index route.
func IndexHandler(indexTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexTemplate.Execute(w, nil)
	}
}

// ViewHandler checks whether the passed ID image exists, and displays the
// view and image.
func ViewHandler(viewTemplate, notFoundTemplate *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		fileInfo, err := goimage.GetFileInfo(config.ImageDirectory, id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Error getting file info.")
			return
		}

		if fileInfo == nil {
			w.WriteHeader(http.StatusBadRequest)
			notFoundTemplate.Execute(w, nil)
			return
		}

		viewTemplate.Execute(w, struct {
			ID       string
			Ext      string
			Filename string
		}{
			ID:       id,
			Ext:      fileInfo.Extension,
			Filename: fileInfo.Filename,
		})
	}
}

// UploadHandler handles file uploads. It takes the POST request, performs
// checks and stores the image.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
}
