package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"metrics/internal/security"
	"net"
	"net/http"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"go.uber.org/zap"
)

type Client struct {
	RestyClient *resty.Client
	Logger      *zap.SugaredLogger
	PublicKey   *rsa.PublicKey
	Key         string
}

func NewClient(serverAddr, key string, cryptoPath string, logger *zap.SugaredLogger) *Client {
	handlerLogger := logger.With("agent", "client")
	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 3
	retryClient.Backoff = customBackoff
	retryClient.CheckRetry = createRetryPolicy(handlerLogger)

	restyClient := resty.NewWithClient(retryClient.StandardClient()).
		SetBaseURL(serverAddr).
		SetHeader("Content-Encoding", "gzip").
		SetHeader("Content-Type", "application/json")

	var publicKey *rsa.PublicKey
	if cryptoPath != "" {
		var err error
		publicKey, err = security.LoadRSAPublicKeyFromCert(cryptoPath)
		if err != nil {
			log.Fatalf("failed to load public key: %v", err)
		}
	}

	client := &Client{
		RestyClient: restyClient,
		Logger:      logger,
		Key:         key,
		PublicKey:   publicKey,
	}

	client.RestyClient.OnBeforeRequest(func(c *resty.Client, r *resty.Request) error {
		return client.processRequest(r)
	})

	return client
}

func (c *Client) processRequest(r *resty.Request) error {
	body := r.Body
	if body == nil {
		return nil
	}

	requestBody, err := readBody(body)
	if err != nil {
		return fmt.Errorf("failed to read request body: %w", err)
	}
	if c.Key != "" {
		hashHex := c.hash(requestBody)
		r.SetHeader("HashSHA256", hashHex)
	}

	if c.PublicKey != nil {
		requestBody, err = c.encrypt(requestBody)
		if err != nil {
			return fmt.Errorf("failed to encrypt request body: %w", err)
		}
	}

	ip, err := getLocalIP()
	if err == nil {
		r.SetHeader("X-Real-IP", ip)
	} else {
		c.Logger.Infoln("unable to detect local IP", "error", err)
	}

	compressedData, err := c.compress(requestBody)
	if err != nil {
		return fmt.Errorf("failed to compress request body: %w", err)
	}
	r.SetBody(compressedData)

	return nil
}

func readBody(body interface{}) ([]byte, error) {
	switch b := body.(type) {
	case []byte:
		return b, nil
	case string:
		return []byte(b), nil
	default:
		return nil, fmt.Errorf("unsupported body type: %T", body)
	}
}

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

	if err = gzipWriter.Close(); err != nil {
		return nil, fmt.Errorf(errClosingGzip, err)
	}

	return buf.Bytes(), nil
}

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

func createRetryPolicy(logger *zap.SugaredLogger) retryablehttp.CheckRetry {
	return func(ctx context.Context, resp *http.Response, err error) (bool, error) {
		if err != nil {
			if resp != nil && resp.StatusCode >= http.StatusInternalServerError {
				logger.Infoln(
					"status_code", resp.StatusCode,
					"retryable", true,
				)
				return true, nil
			}

			var DNSError *net.DNSError
			retry := errors.Is(err, syscall.ECONNREFUSED) ||
				errors.Is(err, syscall.ETIMEDOUT) ||
				err.Error() == "EOF" ||
				errors.As(err, &DNSError)

			if retry {
				logger.Infoln("Connect problem", "retryable", true)
				return true, nil
			}

			if resp == nil {
				logger.Infoln("Connect problem and response == nil and EOF", "retryable", true)
				return true, nil
			}
		}

		return false, nil
	}
}

func customBackoff(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
	backoffIntervals := []time.Duration{1 * time.Second, 3 * time.Second, 5 * time.Second}
	if attemptNum-1 >= 0 && attemptNum-1 < len(backoffIntervals) {
		return backoffIntervals[attemptNum-1]
	}
	return max
}

func (c *Client) hash(data []byte) string {
	return security.HMACSHA256Base64(data, []byte(c.Key))
}

func getLocalIP() (string, error) {
	conns, err := net.Interfaces()
	if err != nil {
		return "", fmt.Errorf("Get local ip: %w", err)
	}
	for _, conn := range conns {
		addrs, err := conn.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			if ip.To4() != nil {
				return ip.String(), nil
			}
		}
	}
	return "", errors.New("cannot find non-loopback IP address")
}
