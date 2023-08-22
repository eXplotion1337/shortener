package internal

import (
	"crypto/rand"
	"math/big"
)

const charset = "0123456789"

func GenerateRandomString(length int) string {
	result := make([]byte, length)
	charsetLength := big.NewInt(int64(len(charset)))

	for i := 0; i < length; i++ {
		n, _ := rand.Int(rand.Reader, charsetLength)
		result[i] = charset[n.Int64()]
	}

	return string(result)
}
