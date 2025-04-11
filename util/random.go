package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// generates random integer between min and max
func RandomInt(min int64, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// generates random string based on length n
func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// generate random name
func GenerateRandomName() string {
	return RandomString(8)
}

// generate random money
func GenerateRandomMoney() int64 {
	return RandomInt(0, 1000)
}

// generate random currency
func GenerateRandomCurrency() string {
	currencies := []string{EUR, USD, CAD}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// generate random email
func GenerateRandomEmail() string {
	return fmt.Sprintf("%s@gmail.com", RandomString(6))
}
