// random
package random

import (
	"math/rand"
	"time"
)

const letterBytes = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "1234567890"

func RandomString(sizeOfRandom int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, sizeOfRandom)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func RandomNumberString(sizeOfRandom int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, sizeOfRandom)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}
