package utils

import (
	"math/rand"
	"time"
)

const alphanumericCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand = rand.New(rand.NewSource(time.Now().UnixNano()))

// GenerateID generates a random alphanumeric string of a specified length.
func GenerateID(length int) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = alphanumericCharset[seededRand.Intn(len(alphanumericCharset))]
	}
	return string(b)
}

// GenerateEmployeeID generates a 10-character alphanumeric ID for an employee.
func GenerateEmployeeID() string {
	return GenerateID(10)
}
