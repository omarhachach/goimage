package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"io/ioutil"
	"encoding/json"
	"strconv"

	"github.com/Omar-H/goimage/util"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type ViewData struct {
	Id string
}

type Config struct {
	Port int `json:"port"`
	Secure bool `json:"secure"`
	AuthKey string `json:"32-byte-auth-key"`
	AllowedMimeTypes []string `json:"allowed-mime-types"`
	AllowedExtensions []string `json:"allowed-extensions"`
	MaxFileSize int64 `json:"max-file-size"`
	ImageDirectory string `json:"image-directory"`
	TemplateDirectory string `json:"template-directory"`
}

var config Config

func main() {
	jsonFiles, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(jsonFiles, &config)

	CSRF := csrf.Protect(
		[]byte(config.AuthKey),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("_csrf"),
		csrf.Secure(config.Secure),
	)
	r := mux.NewRouter()
	templates := template.Must(template.ParseGlob(config.TemplateDirectory + "*.html"))

	r.HandleFunc("/", HomeHandler(templates)).Methods("GET")
	r.HandleFunc("/{id}/", ViewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", UploadHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(config.Port), CSRF(r)))
}

func HomeHandler(t *template.Template) http.HandlerFunc {
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

func ViewHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := ""
		if vars["id"] != "" {
			id = vars["id"]
		}
		err := t.ExecuteTemplate(w, "view.html", ViewData{
			Id: id,
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Fatal(err)
			return
		}
	}
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(config.MaxFileSize)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	if CheckFileType(handler.Header) == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := GenerateName(5)
	f, err := os.OpenFile(config.ImageDirectory + name + util.GetFileExt(handler.Header["Content-Disposition"][0]), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	http.Redirect(w, r, "/" + name + "/", http.StatusSeeOther)
}



func CheckFileType(f textproto.MIMEHeader) bool {
	ext := util.GetFileExt(f["Content-Disposition"][0])

	if !util.Contains(config.AllowedMimeTypes, f["Content-Type"][0]) || !util.Contains(config.AllowedExtensions, ext) {
		return false
	}

	return true
}

func GenerateName(n int) string {
	name := util.GenerateName(n)

	for util.CheckExists(name, "images/") {
		name = util.GenerateName(n)
	}

	return name
}
