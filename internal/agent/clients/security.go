package clients

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func (c *Client) hash(data []byte) string {
	h := hmac.New(sha256.New, []byte(c.key))
	h.Write(data)
	hash := h.Sum(nil)
	return hex.EncodeToString(hash)
}
