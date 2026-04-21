package utils

import (
	"crypto/sha256"
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
