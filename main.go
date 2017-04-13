package main

import (
	"html/template"
	"net/http"
	"log"

	"github.com/gorilla/mux"
)

type Data struct {
	ID string
}

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
		vars := mux.Vars(r)
		id := ""
		if vars["id"] != "" {
			id = vars["id"]
		}
		err := t.ExecuteTemplate(w, "view.html", &Data{ID: id})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	r := mux.NewRouter()
	indexTemplate := template.Must(template.ParseFiles("templates/index.html"))
	viewTemplate := template.Must(template.ParseFiles("templates/view.html"))

	r.HandleFunc("/", homeHandler(indexTemplate)).Methods("GET")
	r.HandleFunc("/{id}/", viewHandler(viewTemplate)).Methods("GET")
	r.HandleFunc("/upload/", uploadHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", r))
}
