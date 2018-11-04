package api

import (
	"net/http"
	"os"

	"github.com/omar-h/goimage"
	"github.com/omar-h/goimage/api/config"
	"github.com/omar-h/goimage/utils"
	"github.com/sirupsen/logrus"
)

// UploadSuccess is the response type for a succesful upload.
type UploadSuccess struct {
	Success bool          `json:"success,omitempty"`
	Code    int           `json:"code,omitempty"`
	File    *goimage.File `json:"file,omitempty"`
}

// UploadHandler handles the "/upload" route.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(config.Cfg.FileBufferSize)

	f, handler, err := r.FormFile("image")
	if err != nil {
		if err == http.ErrMissingFile {
			RenderError(w, ErrMissingFile)
		} else {
			RenderError(w, ErrBadRequest)
		}
		return
	}

	file := goimage.NewFile(f, handler)
	defer file.Close()

	if file.Size > config.Cfg.MaxFileSize {
		RenderError(w, ErrFileTooLarge)
		return
	}

	if !utils.ContainsString(file.Extension, config.Cfg.AllowedExtensions) {
		RenderError(w, ErrFileType)
		return
	}

	err = file.GenerateName(config.Cfg.NameLength).Place(config.Cfg.FileUploadLocation)
	for err != nil {
		if err == os.ErrExist {
			err = file.GenerateName(config.Cfg.NameLength).Place(config.Cfg.FileUploadLocation)
		} else {
			RenderError(w, ErrInternal)
			logrus.WithError(err).Error("An error occured while placing the image.")
			return
		}
	}

	RenderSuccess(w, file)
}
