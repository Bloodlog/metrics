package clients

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {
	logger, _ := zap.NewDevelopment()
	sugar := logger.Sugar()

	serverAddr := "http://example.com"
	key := "test-key"

	client := NewClient(serverAddr, key, sugar)

	assert.NotNil(t, client)

	assert.NotNil(t, client.RestyClient)

	assert.Equal(t, sugar, client.Logger)

	assert.Equal(t, key, client.Key)

	header := client.RestyClient.Header
	assert.Equal(t, serverAddr, client.RestyClient.BaseURL)
	assert.Contains(t, header.Get("Content-Encoding"), "gzip")
	assert.Contains(t, header.Get("Content-Type"), "application/json")
}

func TestReadBody(t *testing.T) {
	tests := []struct {
		name    string
		input   interface{}
		want    []byte
		wantErr bool
	}{
		{
			name:    "Byte slice input",
			input:   []byte("test"),
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name:    "String input",
			input:   "test",
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name:    "Unsupported type (int)",
			input:   42,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Unsupported type (struct)",
			input:   struct{}{},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readBody(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("readBody() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !bytes.Equal(got, tt.want) {
				t.Errorf("readBody() = %v, want %v", got, tt.want)
			}
		})
	}
}
