package clients

import (
	"crypto/rand"
	"crypto/rsa"
	"errors"
	"fmt"
)

func (c *Client) encrypt(data []byte) ([]byte, error) {
	if c.PublicKey == nil {
		return nil, errors.New("public key is not loaded")
	}
	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, c.PublicKey, data)
	if err != nil {
		return nil, fmt.Errorf("encryption failed: %w", err)
	}
	return encryptedData, nil
}
