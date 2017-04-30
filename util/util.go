package util

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

// Variables used for the GenerateName function
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

// Checks if a string slice contains a string
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Checks if the filename already exists in
// a directory.
//
// NOTE: THIS IS GLOB-LIKE AND WILL CHECK IF
// THE DIRECTORY CONTAINS A FILE, WHICH CONTAINS
// THE "name" STRING.
// IT IS EXPLICITLY MADE TO BE USED BY THE
// goimage PROGRAM
func CheckExists(name string, dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.Contains(f.Name(), name) {
			return true
		}
	}
	return false
}

// Generates a name of a specified length
func GenerateName(n int) string {
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

// Gets the file extension from the specified string
//
// NOTE: CHECKS A VERY SPECIFIC STRING (http.Request
// FileForm handler.Header["Content-Disposition"][0])
// AND IT IS EXPLICITLY USED TO BE USED BY THE
// goimage PROGRAM
func GetFileExt(s string) string {
	Ext := strings.Split(s, ";")
	var ext string

	for _, e := range Ext {
		if strings.Contains(e, "filename") {
			ext = filepath.Ext(strings.Trim(strings.Split(e, "=\"")[1], "\""))
		}
	}

	return ext
}
