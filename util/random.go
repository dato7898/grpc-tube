package util

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

func init() {
	rand.New(rand.NewSource(time.Now().UnixNano()))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomCode generate a random n digit code
func RandomCode(n int) string {
	max := int64(math.Pow10(n) - 1)
	code := RandomInt(0, max)

	return fmt.Sprintf("%0*d", n, code)
}
