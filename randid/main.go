// randid Package will generate short encoded unique ids.
package randid

import (
	"math/rand"
	"time"
	"strconv"
)

// RandomId returns a base n (10, 8, 16, 32) encoded random int64 number
func Generate(n int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomId := strconv.FormatInt(r.Int63(), n)

	return randomId
}
