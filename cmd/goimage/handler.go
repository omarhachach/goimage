package main

import (
	"html/template"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/omar-h/goimage"
	"github.com/omar-h/goimage/utils"
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
		err := indexTemplate.Execute(w, nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Error executing template.")
			return
		}
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

		err = viewTemplate.Execute(w, fileInfo)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			logrus.WithError(err).Error("Error executing template.")
			return
		}
	}
}

// UploadHandler handles file uploads. It takes the POST request, performs
// checks and stores the image.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(32 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		logrus.WithError(err).Error("Error handling image.")
		return
	}

	defer file.Close()

	if handler.Size > config.MaxFileSize {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Debug("File is way too big.")
		return
	}

	ext := goimage.GetFileExtension(handler.Filename)
	if !utils.ContainsString(ext, config.AllowedExtensions) {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Debug("The extension isn't allowed.")
		return
	}

	if !utils.ContainsString(goimage.GetFileMIMEType(handler.Header), config.AllowedMIMETypes) {
		w.WriteHeader(http.StatusBadRequest)
		logrus.Debug("The mime type isn't allowed.")
		return
	}

	var id string
	for id == "" {
		id = utils.GenerateName(config.ImageNameLength)

		fileInfo, err := goimage.GetFileInfo(config.ImageDirectory, id)
		if err != nil {
			logrus.WithError(err).Error("Error getting file info.")
			return
		}

		if fileInfo != nil {
			id = ""
		}
	}

	goimage.MoveFile(file, config.ImageDirectory+id+"."+ext)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		logrus.WithError(err).Error("Error moving file.")
		return
	}

	http.Redirect(w, r, "/"+id+"/", http.StatusSeeOther)
}
