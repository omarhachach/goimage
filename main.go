package main

import (
	"html/template"
	"net/http"
	"net/textproto"
	"log"
	"os"
	"io"

	"github.com/gorilla/mux"
	"github.com/gorilla/csrf"
)

type ViewData struct {
	ID string
}

func homeHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-CSRF-Token", csrf.Token(r))
		err := t.ExecuteTemplate(w, "index.html", map[string]interface{}{
			csrf.TemplateTag: csrf.TemplateField(r),
		})
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
		err := t.ExecuteTemplate(w, "view.html", ViewData{
			ID: id,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(30000000)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()
	log.Printf("%d", handler.Header)

	checkFileType(handler.Header)
	name := generateName()
	f, err := os.OpenFile("./images/" + name, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
}

func checkFileType(f textproto.MIMEHeader) string, error {

}

func generateName() string {
	return "123"
}

func main() {
	CSRF := csrf.Protect(
		[]byte("62caed6a7842b5470c2e89693f92c9bab01219f8ebc0c9c0785b97cfd7a68187"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("_csrf"),
		csrf.Secure(false),
	)
	r := mux.NewRouter()
	templates := template.Must(template.ParseGlob("templates/*.html"))

	r.HandleFunc("/", homeHandler(templates)).Methods("GET")
	r.HandleFunc("/{id}/", viewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", uploadHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", CSRF(r)))
}
