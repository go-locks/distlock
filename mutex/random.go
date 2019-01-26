package mutex

import (
	"math/rand"
	"time"
)

var (
	randLen      = 1 << 5
	byteSlice    = []byte("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	byteSliceLen = len(byteSlice)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func generateValue() string {
	result := make([]byte, randLen)
	for i := 0; i < randLen; i++ {
		result[i] = byteSlice[rand.Intn(byteSliceLen)]
	}
	return string(result)
}
