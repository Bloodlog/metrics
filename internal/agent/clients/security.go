package clients

import (
	"metrics/internal/security"
)

func (c *Client) hash(data []byte) string {
	return security.HMACSHA256Base64(data, []byte(c.Key))
}
