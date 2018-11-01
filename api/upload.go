package api

import (
	"net/http"
	"os"

	"github.com/omar-h/goimage"
	"github.com/omar-h/goimage/utils"
)

var allowedExts = []string{
	".png",
	".jpg",
	".jpeg",
	".jiff",
	".ico",
	".gif",
	".tif",
	".webp",
}

const maxFileSize = 20 << 18 // 5 MB

// UploadHandler handles the "/upload" route.
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(maxFileSize)

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

	if file.Size > maxFileSize {
		RenderError(w, ErrFileTooLarge)
		return
	}

	if !utils.ContainsString(file.Extension, allowedExts) {
		RenderError(w, ErrFileType)
		return
	}

	err = file.GenerateName(6).Place("img/")
	for err != nil {
		if err == os.ErrExist {
			err = file.GenerateName(6).Place("img/")
		} else {
			RenderError(w, ErrInternal)
			return
		}
	}

	http.Redirect(w, r, "/"+file.Basename+"/", http.StatusSeeOther)
}
