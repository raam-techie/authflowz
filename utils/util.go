package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// GenerateAccountID creates an 11-digit account ID starting with "112".
// Format: 112XXXXXXXX  (3 fixed prefix + 8 random digits)
func GenerateAccountID() (string, error) {
	const prefix = "112"
	const randomDigits = 8

	result := prefix
	for i := 0; i < randomDigits; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", fmt.Errorf("failed to generate account id: %w", err)
		}
		result += n.String()
	}
	return result, nil
}
