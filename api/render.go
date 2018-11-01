package api

import (
	"net/http"

	"github.com/unrolled/render"
)

// Render is the global renderer for the API.
var Render = render.New()

// RenderError is a helper type that will render a JSON error.
func RenderError(w http.ResponseWriter, err ErrorResponse) {
	Render.JSON(w, err.Code, err)
}
