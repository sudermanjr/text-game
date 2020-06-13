package utils

import (
	"math/rand"
)

// RandomBool returns a random true/false
func RandomBool() bool {
	return rand.Int31()&0x01 == 0
}
