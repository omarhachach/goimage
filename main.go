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

// Data parsed to the ViewHandler's template
type ViewData struct {
	Id       string
	ImageUrl string
	Ext      string
}

// Parsed JSON config
type Config struct {
	Port              int      `json:"port"`
	Secure            bool     `json:"secure"`
	AuthKey           string   `json:"32-byte-auth-key"`
	AllowedMimeTypes  []string `json:"allowed-mime-types"`
	AllowedExtensions []string `json:"allowed-extensions"`
	ImageNameLength   int      `json:"image-name-length"`
	MaxFileSize       int64    `json:"max-file-size"`
	ImageDirectory    string   `json:"image-directory"`
	TemplateDirectory string   `json:"template-directory"`
	PublicDirectory   string   `json:"public-directory"`
	ImageUrl          string   `json:"image-url"`
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
	err = os.MkdirAll(config.TemplateDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Template Directory")
	}
	err = os.MkdirAll(config.PublicDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Public Directory")
	}
	err = os.MkdirAll(config.ImageDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Image Directory")
	}

	templates := template.Must(template.ParseGlob(config.TemplateDirectory + "*.html"))

	r.HandleFunc("/{id}/", ViewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", UploadHandler).Methods("POST")
	r.PathPrefix("/").HandlerFunc(RootHandler(templates)).Methods("GET")

	log.Print("Listening on port " + strconv.Itoa(config.Port))
	if config.CSRF {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), CSRF(r)))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(config.Port), r))
	}
}

// RootHandler handles the root route.
// This includes the homepage and the
// file system.
func RootHandler(t *template.Template) http.HandlerFunc {
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

// Handles the ability to view your image.
// It checks if the file exists, and returns
// an appropriate response.
//
// If the image exists, it serves a template
// which is defined in
// TemplateDirectory/view.html.
func ViewHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		ext := util.GetFileExtFromDir(id, config.ImageDirectory)
		if ext != "" {
			err := t.ExecuteTemplate(w, "view.html", ViewData{
				Id:       id,
				ImageUrl: config.ImageUrl,
				Ext:      ext,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
		} else {
			http.ServeFile(w, r, config.TemplateDirectory+"404.html")
		}
	}
}

// Handles the upload route.
// It checks the file, and uploads the
// image to the ImageDirectory defined
// in the config.json
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

	name := GenerateName(config.ImageNameLength)
	f, err := os.OpenFile(config.ImageDirectory+name+util.GetFileExt(handler.Header["Content-Disposition"][0]), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	http.Redirect(w, r, "/"+name+"/", http.StatusSeeOther)
}

// CheckFileType checks if the uploaded
// file's extension is allowed.
// It accepts a textproto.MIMEHeader
// which is recieved from the http.Request
// FormFile function, as the handler.Header
func CheckFileType(f textproto.MIMEHeader) bool {
	ext := util.GetFileExt(f["Content-Disposition"][0])

	if !util.Contains(config.AllowedMimeTypes, f["Content-Type"][0]) || !util.Contains(config.AllowedExtensions, ext) {
		return false
	}

	return true
}

// GenerateName generates a name with
// a specified length, and checks if
// the name already exists, using the
// util functions for this.
func GenerateName(n int) string {
	name := util.GenerateName(n)

	for util.CheckExists(name, config.ImageDirectory) {
		name = util.GenerateName(n)
	}

	return name
}
