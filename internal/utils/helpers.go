package utils

import (
	"math/rand"
	"time"
)

func GenerateVerificationCode() string {
	const charset = "0123456789"
	rand.Seed(time.Now().UnixNano())
	code := make([]byte, 4)
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}
