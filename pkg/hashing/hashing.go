package hashing

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func generateRandom(size int) ([]byte, error) {
	b := make([]byte, size)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func GenerateHash(key string, data []byte) (string, error) {
	secretKey := sha256.Sum256([]byte(key))
	aesblock, err := aes.NewCipher(secretKey[:])
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return "", err
	}

	nonce, err := generateRandom(aesgcm.NonceSize())
	if err != nil {
		return "", err
	}

	dst := aesgcm.Seal(nonce, nonce, data, nil)
	encryptedStr := fmt.Sprintf("%x", dst)
	return encryptedStr, nil
}

func VerifyHash(key string, data []byte, encryptedHash string) (bool, error) {
	encryptedHashBytes, err := hex.DecodeString(encryptedHash)
	if err != nil {
		return false, err
	}

	secretKey := sha256.Sum256([]byte(key))
	aesblock, err := aes.NewCipher(secretKey[:])
	if err != nil {
		return false, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return false, err
	}

	nonceSize := aesgcm.NonceSize()
	nonce := encryptedHashBytes[:nonceSize]
	encryptedData := encryptedHashBytes[nonceSize:]

	decryptedData, err := aesgcm.Open(nil, nonce, encryptedData, nil)
	if err != nil {
		return false, err
	}
	return compareSlices(decryptedData, data), nil
}

func compareSlices(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
