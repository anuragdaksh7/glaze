package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"
)

func GenerateHashID() string {
	seed := fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1e6))
	hash := sha256.Sum256([]byte(seed))
	return hex.EncodeToString(hash[:])[:16]
}

func GenerateRandomString(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		// In a real app, handle this more gracefully, but
		// if crypto/rand fails, the system has bigger problems.
		panic(fmt.Sprintf("failed to generate random bytes: %v", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}
