package utils

import (
	"math/rand"
	"time"
)

// Variables used for the GenerateName function
const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVXYZ1234567890"
	letterIdxBits = 6
	letterIdxMask = 1<<letterIdxBits - 1
	letterIdxMax  = 63 / letterIdxBits
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

// GenerateName generates a name of the specified length.
func GenerateName(length int) (name string) {
	var src = rand.NewSource(time.Now().UnixNano())

	b := make([]byte, length)
	for i, cache, remain := length-1, src.Int63(), letterIdxMax; i >= 0; {
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
