package main

import (
	"html/template"
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

func homeHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := t.ExecuteTemplate(w, "index.html", nil)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func viewHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := ""
		err := t.ExecuteTemplate(w, "index.html", id)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func main() {
	r := mux.NewRouter()
	indexTemplate := template.Must(template.ParseFiles("templates/index.html"))

	r.HandleFunc("/", homeHandler(indexTemplate)).Methods("GET")
	r.HandleFunc("/{id}/", viewHandler)

	log.Fatal(http.ListenAndServe(":8080", r))
}
