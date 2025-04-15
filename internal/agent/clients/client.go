package clients

import (
	"crypto/rsa"
	"fmt"
	"log"
	"metrics/internal/security"

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
