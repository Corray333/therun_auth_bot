package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/nyaruka/phonenumbers"
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

func FormatPhoneNumber(userPhone string, savedPhone string) (string, error) {
	// Parse the user's phone number
	parsedNumber := &phonenumbers.PhoneNumber{}
	var err error
	f := true
	for _, reg := range codes[strings.Split(savedPhone, " ")[0]] {
		parsedNumber, err = phonenumbers.Parse(userPhone, reg)
		if err != nil {
			continue
		}
		f = false
		break
	}
	if f {
		return "", fmt.Errorf("Something went wrong, try to use sms verification.")
	}

	// Check if the number is valid
	if !phonenumbers.IsValidNumber(parsedNumber) {
		return "", fmt.Errorf("invalid phone number")
	}

	// Format the number to E164 format
	formattedNumber := phonenumbers.Format(parsedNumber, phonenumbers.E164)
	return formattedNumber, nil
}
