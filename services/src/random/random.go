// random
package random

import (
	"math/rand"
)

const letterBytes = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const numberBytes = "1234567890"

func RandomString(sizeOfRandom int) string {
	b := make([]byte, sizeOfRandom)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
func RandomNumberString(sizeOfRandom int) string {
	b := make([]byte, sizeOfRandom)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}
	return string(b)
}
