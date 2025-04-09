package clients

import (
	"bytes"
	"compress/gzip"
	"io"
	"testing"
)

func TestClient_Compress(t *testing.T) {
	client := &Client{}

	t.Run("compress valid data", func(t *testing.T) {
		data := []byte("hello world")
		compressed, err := client.compress(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gzipReader, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer func(gzipReader *gzip.Reader) {
			err = gzipReader.Close()
			if err != nil {
				t.Fatalf("failed to create decompressed data: %v", err)
			}
		}(gzipReader)

		decompressed, err := io.ReadAll(gzipReader)
		if err != nil {
			t.Fatalf("failed to read decompressed data: %v", err)
		}

		if !bytes.Equal(decompressed, data) {
			t.Fatalf("expected %s, got %s", data, decompressed)
		}
	})

	t.Run("compress empty data", func(t *testing.T) {
		var data []byte
		compressed, err := client.compress(data)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		gzipReader, err := gzip.NewReader(bytes.NewReader(compressed))
		if err != nil {
			t.Fatalf("failed to create gzip reader: %v", err)
		}
		defer func(gzipReader *gzip.Reader) {
			err = gzipReader.Close()
			if err != nil {
				t.Fatalf("failed to read decompressed data: %v", err)
			}
		}(gzipReader)

		decompressed, err := io.ReadAll(gzipReader)
		if err != nil {
			t.Fatalf("failed to read decompressed data: %v", err)
		}

		if !bytes.Equal(decompressed, data) {
			t.Fatalf("expected empty data, got %s", decompressed)
		}
	})
}
