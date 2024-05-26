package util

import (
	"math/rand"
	"time"
)

func RandomNumberGenerator(limit int) int {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate a random number between 1 and limit
	randomNumber := rand.Intn(limit) + 1

	return randomNumber
}
