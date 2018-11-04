package utils

import (
	"math/rand"
	"strings"
)

// ContainsString checks if a string slice contains a string.
func ContainsString(item string, slice []string) (contains bool) {
	for _, a := range slice {
		if a == item {
			return true
		}
	}

	return false
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// GenerateName generates a name of the specified length.
func GenerateName(strLen int) string {
	b := make([]byte, strLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

// GetFileBasename returns the base file name of a given filename.
// Eg. the file name without the extension.
func GetFileBasename(filename string) (basename string) {
	index := strings.LastIndex(filename, ".")
	if index == -1 {
		return filename
	}

	return filename[:index]
}
