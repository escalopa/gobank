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

// Creates a random integer from min to max
func RandomInteger(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		sb.WriteByte(alphabet[rand.Intn(k)])
	}

	return sb.String()
}

func RandomOwner() string {
	return RandomString(7)
}

func RandomUsername() string {
	return RandomString(10)
}

func RandomMoney() int64 {
	return RandomInteger(0, 1000)
}

func RandomCurrency() string {
	currency := []string{"EGP", "USD", "RUB"}
	return currency[rand.Intn(len(currency))]
}

func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(7))
}
