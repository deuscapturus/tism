// randid Package will generate short encoded unique ids.
package randid

import (
	"math/rand"
	"strconv"
	"time"
)

var myrand *rand.Rand

func init() {
	myrand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomId returns a base n (10, 8, 16, 32) encoded random int64 number
func Generate(n int) string {
	randomId := strconv.FormatInt(myrand.Int63(), n)

	return randomId
}

func GenerateSecret(size int) string {

	const chars = "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz0123456789@$%^&*"
	secret := make([]byte, size)
	for i := 0; i < size; i++ {
		secret[i] = chars[myrand.Intn(len(chars))]
	}

	return string(secret)
}
