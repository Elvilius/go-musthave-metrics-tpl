package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashingTrue(t *testing.T) {
	data := []byte("test data")
	hash := GenerateHash("secret", data)
	ok := VerifyHash("secret", data, hash)
	assert.True(t, ok)
}
