package clients

import (
	"bytes"
	"compress/gzip"
	"fmt"
)

func (c *Client) compress(data []byte) ([]byte, error) {
	const (
		errCompressingData = "error compressing the data: %w"
		errClosingGzip     = "error closing gzip stream: %w"
	)
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write(data)
	if err != nil {
		return nil, fmt.Errorf(errCompressingData, err)
	}

	if err := gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf(errClosingGzip, err)
	}

	return buf.Bytes(), nil
}
