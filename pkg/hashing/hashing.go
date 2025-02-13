// Package hashing provides functions for generating and verifying HMAC SHA-256 hashes.
package hashing

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// GenerateHash creates an HMAC SHA-256 hash of the given data using the provided key.
//
// The result is a base64-encoded string representation of the hash.
//
// Parameters:
//   - key: The secret key used to generate the hash.
//   - data: The data to be hashed.
//
// Returns:
//   - A base64-encoded HMAC SHA-256 hash.
func GenerateHash(key string, data []byte) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)

	sum := h.Sum(nil)
	return base64.StdEncoding.EncodeToString(sum)
}

// VerifyHash compares the provided encrypted hash with the HMAC SHA-256 hash
// generated from the given key and data.
//
// This function ensures data integrity by verifying that the computed hash
// matches the expected value.
//
// Parameters:
//   - key: The secret key used to generate the hash.
//   - data: The original data that was hashed.
//   - encryptedHash: The base64-encoded hash to verify.
//
// Returns:
//   - true if the hash matches, false otherwise.
func VerifyHash(key string, data []byte, encryptedHash string) bool {
	hash := GenerateHash(key, data)
	return hash == encryptedHash
}