package nanoid

import (
	"bytes"
	"math/rand"
	"time"
)

const (
	defaultAlphabet = "0123456789"
	defaultLength   = 6
)

var random = rand.New(rand.NewSource(time.Now().UnixNano()))

// New generates a random nanoid with length of [defaultLength]
func New() string {
	return NewWithLength(defaultLength)
}

// Generate generates a random nanoid with a custom length; 10 <= length <= 200
// If the length was out of that range, it defaults to [defaultLength].
func NewWithLength(length int) string {
	if length < 10 || length > 200 {
		length = defaultLength
	}
	buf := bytes.NewBuffer([]byte{})
	for i := 0; i < length; i++ {
		buf.WriteByte(defaultAlphabet[random.Intn(len(defaultAlphabet))])
		random.Seed(time.Now().UnixNano())
	}
	return buf.String()
}
