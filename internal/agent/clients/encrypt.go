package clients

import (
	"errors"
	"fmt"
	"metrics/internal/security"
)

func (c *Client) encrypt(data []byte) ([]byte, error) {
	if c.PublicKey == nil {
		return nil, errors.New("public key is not loaded")
	}
	encrypted, err := security.EncryptWithPublicKey(data, c.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("encrypt with public key: %w", err)
	}
	return encrypted, nil
}
