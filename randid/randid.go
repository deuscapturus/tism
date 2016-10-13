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
