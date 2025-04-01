package clients

import (
	"testing"
)

func TestHash(t *testing.T) {
	client := &Client{
		Key: "your-secret-key",
	}

	data := []byte("test data")
	expectedHash := "uISqqd+u7OkFYcS8P+wNxry5N9deUZBzocZ8HQMqARA="

	actualHash := client.hash(data)

	if actualHash != expectedHash {
		t.Errorf("expected %s, got %s", expectedHash, actualHash)
	}
}
