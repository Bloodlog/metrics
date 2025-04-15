package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
)

// HMACSHA256Base64 возвращает base64(hmac(sha256(data, key))).
func HMACSHA256Base64(data, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// CheckHMACSHA256Base64 сравнивает base64(hmac(sha256(data, key))) с переданным хешем.
func CheckHMACSHA256Base64(data, key []byte, providedBase64 string) bool {
	expected := HMACSHA256Base64(data, key)
	return hmac.Equal([]byte(expected), []byte(providedBase64))
}
