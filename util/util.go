package util

import (
	"io/ioutil"
	"math/rand"
	"path/filepath"
	"strings"
	"time"
)

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
)

func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func CheckExists(name string, dir string) bool {
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if strings.Contains(f.Name(), name) {
			return true
		}
	}
	return false
}

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