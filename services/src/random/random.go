// random
package random

import (
	"math/rand"
)

const letterBytes = "1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandomString(sizeOfRandom int) string {
	b := make([]byte, sizeOfRandom)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
