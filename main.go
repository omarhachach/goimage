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

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/omar-h/goimage/util"
)

// viewData is the data parsed to the
// ViewHandler's template
type viewData struct {
	ID       string
	ImageURL string
	Ext      string
}

// config is the parsed JSON config
type config struct {
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
	ImageURL          string   `json:"image-url"`
	CSRF              bool     `json:"csrf"`
}

var parsedConfig config

func main() {
	jsonFiles, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
		return
	}
	json.Unmarshal(jsonFiles, &parsedConfig)

	CSRF := csrf.Protect(
		[]byte(parsedConfig.AuthKey),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("_csrf"),
		csrf.Secure(parsedConfig.Secure),
	)
	r := mux.NewRouter()
	err = os.MkdirAll(parsedConfig.TemplateDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Template Directory")
	}
	err = os.MkdirAll(parsedConfig.PublicDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Public Directory")
	}
	err = os.MkdirAll(parsedConfig.ImageDirectory, 644)
	if err != nil {
		log.Fatal("Unable to create Image Directory")
	}

	templates := template.Must(template.ParseGlob(parsedConfig.TemplateDirectory + "*.html"))

	r.HandleFunc("/{id}/", viewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", uploadHandler).Methods("POST")
	r.PathPrefix("/").HandlerFunc(rootHandler(templates)).Methods("GET")

	log.Print("Listening on port " + strconv.Itoa(parsedConfig.Port))
	if parsedConfig.CSRF {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(parsedConfig.Port), CSRF(r)))
	} else {
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(parsedConfig.Port), r))
	}
}

// rootHandler handles the root route.
// This includes the homepage and the
// file system.
func rootHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch url := r.URL.Path; url {
		case "/":
			if parsedConfig.CSRF {
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
			if _, err := os.Stat(parsedConfig.PublicDirectory + url); err != nil {
				w.Header().Add("Content-Type", "text/html; charset=utf-8")
				w.WriteHeader(http.StatusNotFound)
				http.ServeFile(w, r, parsedConfig.TemplateDirectory+"404.html")
			} else {
				http.ServeFile(w, r, parsedConfig.PublicDirectory+url)
			}
		}

	}
}

// viewHandler handles the ability to view your image.
// It checks if the file exists, and returns
// an appropriate response.
//
// If the image exists, it serves a template
// which is defined in
// TemplateDirectory/view.html.
func viewHandler(t *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"]
		ext := util.GetFileExtFromDir(id, parsedConfig.ImageDirectory)
		if ext != "" {
			err := t.ExecuteTemplate(w, "view.html", viewData{
				ID:       id,
				ImageURL: parsedConfig.ImageURL,
				Ext:      ext,
			})
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Fatal(err)
				return
			}
		} else {
			w.Header().Add("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, parsedConfig.TemplateDirectory+"404.html")
		}
	}
}

// uploadHandler handles the upload route.
// It checks the file, and uploads the
// image to the ImageDirectory defined
// in the config.json
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(30000000)
	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Fatal(err)
		return
	}
	defer file.Close()

	if checkFileType(handler.Header) == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := generateName(parsedConfig.ImageNameLength)
	f, err := os.OpenFile(parsedConfig.ImageDirectory+name+util.GetFileExt(handler.Header["Content-Disposition"][0]), os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
		return
	}

	io.Copy(f, file)
	fStat, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}
	size := fStat.Size()
	if size > parsedConfig.MaxFileSize {
		f.Close()
		err = os.RemoveAll(parsedConfig.ImageDirectory + name + util.GetFileExt(handler.Header["Content-Disposition"][0]))
		if err != nil {
			log.Fatal(err)
		}
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/"+name+"/", http.StatusSeeOther)
		f.Close()
	}
}

// checkFileType checks if the uploaded
// file's extension is allowed.
// It accepts a textproto.MIMEHeader
// which is received from the http.Request
// FormFile function, as the handler.Header
func checkFileType(f textproto.MIMEHeader) bool {
	ext := util.GetFileExt(f["Content-Disposition"][0])

	if !util.Contains(parsedConfig.AllowedMimeTypes, f["Content-Type"][0]) || !util.Contains(parsedConfig.AllowedExtensions, ext) {
		return false
	}

	return true
}

// generateName generates a name with
// a specified length, and checks if
// the name already exists, using the
// util functions for this.
func generateName(n int) string {
	name := util.GenerateName(n)

	for util.CheckExists(name, parsedConfig.ImageDirectory) {
		name = util.GenerateName(n)
	}

	return name
}
