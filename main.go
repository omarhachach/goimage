package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strconv"

	"github.com/Omar-H/goimage/util"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

type ViewData struct {
	Id string
}

type Config struct {
	Port              int      `json:"port"`
	Secure            bool     `json:"secure"`
	AuthKey           string   `json:"32-byte-auth-key"`
	AllowedMimeTypes  []string `json:"allowed-mime-types"`
	AllowedExtensions []string `json:"allowed-extensions"`
	MaxFileSize       int64    `json:"max-file-size"`
	ImageDirectory    string   `json:"image-directory"`
	TemplateDirectory string   `json:"template-directory"`
	PublicDirectory   string   `json:"public-directory"`
	CSRF              bool     `json:"csrf"`
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

	os.Mkdir(config.ImageDirectory, 644)
	os.Mkdir(config.TemplateDirectory, 644)

	r.HandleFunc("/{id}/", ViewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", UploadHandler).Methods("POST")
	r.PathPrefix("/").HandlerFunc(HomeHandler(templates)).Methods("GET")

	log.Print("Listening on port: " + strconv.Itoa(config.Port))
	if config.CSRF {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), CSRF(r)))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), r))
	}
}

func HomeHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch url := r.URL.Path; url {
		case "/":
			if config.CSRF {
				w.Header().Set("X-CSRF-Token", csrf.Token(r))
				err := t.ExecuteTemplate(w, "index.html", map[string]interface{}{
					csrf.TemplateTag: csrf.TemplateField(r),
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			} else {
				err := t.ExecuteTemplate(w, "index.html", map[string]interface{}{
					csrf.TemplateTag: "",
				})
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
				}
			}
		default:
			if _, err := os.Stat(config.PublicDirectory + url); err != nil {
				http.ServeFile(w, r, config.TemplateDirectory+"404.html")
			} else {
				http.ServeFile(w, r, config.PublicDirectory+url)
			}
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
	f, err := os.OpenFile(config.ImageDirectory+name+util.GetFileExt(handler.Header["Content-Disposition"][0]), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	http.Redirect(w, r, "/"+name+"/", http.StatusSeeOther)
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
