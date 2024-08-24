package hashing

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateHash(t *testing.T) {

	data := []byte("zalupa")
	hash, _ := GenerateHash("secret", data)
	fmt.Println(hash)
	ok, _ := VerifyHash("secret", data, hash)
	assert.True(t, ok)
}
