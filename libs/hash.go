package libs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"math/big"
	"strings"
)

const base62Chars = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// EncodeBase62 converts a string to a Base62 encoded string
func encodeBase62(input string) string {
	num := new(big.Int).SetBytes([]byte(input))
	base := big.NewInt(int64(len(base62Chars)))
	var encoded strings.Builder

	zero := big.NewInt(0)
	for num.Cmp(zero) > 0 {
		mod := new(big.Int)
		num.DivMod(num, base, mod)
		encoded.WriteByte(base62Chars[mod.Int64()])
	}

	// Handle the special case for empty input
	if encoded.Len() == 0 {
		encoded.WriteByte(base62Chars[0])
	}

	// Reverse the result since we are constructing the string backwards
	return reverseString(encoded.String())
}

// Helper function to reverse a string
func reverseString(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

// GenerateRandomSalt creates a random salt
func generateRandomSalt(length int) (string, error) {
	salt := make([]byte, length)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func ComputeShortHash(url string, db *map[string]string) string {
	salt, _ := generateRandomSalt(4)
	for {
		base62 := encodeBase62(url)
		hash := md5.Sum([]byte(base62))
		hashString := hex.EncodeToString(hash[:]) // Convert full hash to string
		shortHash := hashString[:7]

		if _, ok := (*db)[shortHash]; !ok {
			return shortHash
		}
		url = url + salt
	}
}
