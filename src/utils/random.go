package utils

import (
	"math/rand"
	"strings"
)

// Returns an int >= min, < max
func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

// Generate a random string of A-Z chars with len = l
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

func GetFinalPathString(url string, val int) string{
	path := strings.Split(url, "/")
	return path[len(path) - val]
}
