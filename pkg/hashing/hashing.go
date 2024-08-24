package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

func GenerateHash(key string, data []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

func VerifyHash(key string, data []byte, encryptedHash string) bool {
	hash := GenerateHash(key, data)
	return hash == encryptedHash
}
