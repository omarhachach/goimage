package main

import (
	"net/http"
)

// MaxBodySizeMiddleware ensures that the request body is of a certain size.
func MaxBodySizeMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" {
			r.Body = http.MaxBytesReader(w, r.Body, config.MaxFileSize)
		}

		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
