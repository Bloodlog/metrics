package security

import "testing"

func TestHMACSHA256Base64(t *testing.T) {
	data := []byte("test-data")
	key := []byte("secret-key")

	expected := HMACSHA256Base64(data, key)
	if expected == "" {
		t.Fatal("expected non-empty hash")
	}

	if !CheckHMACSHA256Base64(data, key, expected) {
		t.Errorf("hash check failed: expected true for valid hash")
	}

	invalidHash := "invalidhash"
	if CheckHMACSHA256Base64(data, key, invalidHash) {
		t.Errorf("hash check passed with invalid hash, expected false")
	}

	wrongKey := []byte("wrong-key")
	if CheckHMACSHA256Base64(data, wrongKey, expected) {
		t.Errorf("hash check passed with wrong key, expected false")
	}
}
