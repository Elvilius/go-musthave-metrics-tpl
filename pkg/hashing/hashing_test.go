package hashing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashingTrue(t *testing.T) {
	data := []byte("test data")
	hash, err := GenerateHash("secret", data)
	assert.NoError(t, err)
	ok, err := VerifyHash("secret", data, hash)
	assert.NoError(t, err)
	assert.True(t, ok)
}

func TestHashingFalse(t *testing.T) {
	data := []byte("test data")
	hash, err := GenerateHash("secret", data)
	assert.NoError(t, err)
	ok, err := VerifyHash("secrets", data, hash)
	assert.Error(t, err)
	assert.False(t, ok)
}
