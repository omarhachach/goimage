package api

import (
	"net/http"

	"github.com/omar-h/goimage"
	"github.com/unrolled/render"
)

// Render is the global renderer for the API.
var Render = render.New()

// RenderError is a helper type that will render a JSON error.
func RenderError(w http.ResponseWriter, err ErrorResponse) {
	Render.JSON(w, err.Code, err)
}

// RenderSuccess will render a success response.
func RenderSuccess(w http.ResponseWriter, file *goimage.File) {
	Render.JSON(w, http.StatusOK, UploadSuccess{
		Success: true,
		Code:    http.StatusOK,
		File:    file,
	})
}
