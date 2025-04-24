package test

import (
	"testing"

	"github.com/Elvilius/go-musthave-metrics-tpl/pkg/crypto"
	"github.com/stretchr/testify/assert"
)

func TestCrypto_Encrypt(t *testing.T) {
	type args struct {
		message []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "Success",
			args: args{
				message: []byte("test"),
			},
			want: []byte("test"),
		},
	}

	cfg := crypto.Cfg{
		PublicKeyPath:  "./public.pem",
		PrivateKeyPath: "./private.pem",
	}

	c, err := crypto.New(cfg)
	assert.Equal(t, err, nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encryptData, err := c.Encrypt(tt.args.message)
			assert.Equal(t, err, nil)

			decryptData, err := c.Decrypt(encryptData)
			assert.Equal(t, err, nil)

			assert.Equal(t, decryptData, tt.want)
		})
	}
}
