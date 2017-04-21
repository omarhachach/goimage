package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"path/filepath"
	"math/rand"
	"time"
	"io/ioutil"
	"encoding/json"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
)

const (
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

type ViewData struct {
	Id string
}

type Config struct {
	Port int `json:"port"`
	Secure bool `json:"secure"`
	AuthKey string `json:"32-byte-auth-key"`
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
			Id: id,
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

	if checkFileType(handler.Header) == false {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	name := generateName(5)
	f, err := os.OpenFile("./images/" + name + getFileExt(handler.Header["Content-Disposition"][0]), os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()

	io.Copy(f, file)
	http.Redirect(w, r, "/" + name + "/", http.StatusSeeOther)
}

func checkFileType(f textproto.MIMEHeader) bool {
	allowed := []string{
		"image/x-icon",
		"image/jpeg",
		"image/pjpeg",
		"image/png",
		"image/tiff",
		"image/x-tiff",
		"image/webp",
		"image/gif",
	}

	allowedExt := []string{
		".png",
		".jpeg",
		".jpg",
		".jiff",
		".png",
		".ico",
		".gif",
		".tif",
		".webp",
	}

	ext := getFileExt(f["Content-Disposition"][0])

	if !contains(allowed, f["Content-Type"][0]) || !contains(allowedExt, ext) {
		return false
	}

	return true
}

func getFileExt(s string) string {
	Ext := strings.Split(s, ";")
	var ext string

	for _, e := range Ext {
		if strings.Contains(e, "filename") {
			ext = filepath.Ext(strings.Trim(strings.Split(e, "=\"")[1], "\""))
		}
	}

	return ext
}

func generateName(n int) string {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, n)
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func main() {
	jsonFiles, err := ioutil.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
		return
	}

	var config Config
	json.Unmarshal(jsonFiles, &config)

	CSRF := csrf.Protect(
		[]byte(config.AuthKey),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.FieldName("_csrf"),
		csrf.Secure(config.Secure),
	)
	r := mux.NewRouter()
	templates := template.Must(template.ParseGlob("templates/*.html"))

	r.HandleFunc("/", homeHandler(templates)).Methods("GET")
	r.HandleFunc("/{id}/", viewHandler(templates)).Methods("GET")
	r.HandleFunc("/upload/", uploadHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(":" + strconv.Itoa(config.Port), CSRF(r)))
}
