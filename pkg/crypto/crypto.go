package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

type Crypto struct {
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
}

type Cfg struct {
	PublicKeyPath  string
	PrivateKeyPath string
}

func New(cfg Cfg) (*Crypto, error) {
	var publicKey *rsa.PublicKey
	var privateKey *rsa.PrivateKey
	var err error

	if cfg.PrivateKeyPath != "" {
		privateKey, err = loadPrivateKey(cfg.PrivateKeyPath)
		if err != nil {
			return nil, err
		}
	}

	if cfg.PublicKeyPath != "" {
		publicKey, err = loadPublicKey(cfg.PublicKeyPath)
		if err != nil {
			return nil, err
		}
	}
	return &Crypto{
		privateKey: privateKey,
		publicKey:  publicKey,
	}, nil
}

func loadPrivateKey(path string) (*rsa.PrivateKey, error) {
	var privateKey *rsa.PrivateKey

	block, err := loadPemData(path)
	if err != nil {
		return nil, err
	}

	privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	if err == nil {
		return privateKey, nil
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey = key.(*rsa.PrivateKey)
	return privateKey, nil
}

func loadPublicKey(path string) (*rsa.PublicKey, error) {
	var publicKey *rsa.PublicKey

	block, err := loadPemData(path)
	if err != nil {
		return nil, err
	}

	publicKey, err = x509.ParsePKCS1PublicKey(block.Bytes)
	if err == nil {
		return publicKey, nil
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey = key.(*rsa.PublicKey)
	return publicKey, nil
}

func loadPemData(path string) (*pem.Block, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to parse PEM block")

	}
	return block, nil
}

func (c *Crypto) Encrypt(message []byte) ([]byte, error) {
	if c.publicKey == nil {
		return nil, errors.New("public key is not loaded")
	}
	return rsa.EncryptPKCS1v15(rand.Reader, c.publicKey, message)
}

func (c *Crypto) Decrypt(ciphertext []byte) ([]byte, error) {
	if c.privateKey == nil {
		return nil, errors.New("private key is not loaded")
	}
	return rsa.DecryptPKCS1v15(rand.Reader, c.privateKey, ciphertext)
}
