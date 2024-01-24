package utils

import (
	"fmt"
	"math/rand"
	"strings"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandomInt generate a random value between the min and max arguments
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomEmail generate a random email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}

// RandomString generate a random string according the size from the arguments
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)
	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]

		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomOwner generate a new random owner name
func RandomOwner() string {
	return RandomString(6)
}

// RandomMoney generate a new random money
func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// RandomCurrency return a random currency
func RandomCurrency() string {
	currencies := []string{"USD", "EUR", "CAD"}

	return currencies[rand.Intn(len(currencies))]
}
