package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
)

// LoadRSAPublicKeyFromCert загружает RSA-публичный ключ из сертификата (PEM-формат).
func LoadRSAPublicKeyFromCert(path string) (*rsa.PublicKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read cert file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("failed to decode PEM block from cert file")
	}

	if block.Type != "CERTIFICATE" {
		return nil, fmt.Errorf("expected CERTIFICATE block, got: %s", block.Type)
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse x509 certificate: %w", err)
	}

	pubKey, ok := cert.PublicKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("certificate does not contain RSA public key")
	}

	return pubKey, nil
}

// LoadRSAPrivateKeyFromFile загружает RSA-приватный ключ из сертификата (PEM-формат).
func LoadRSAPrivateKeyFromFile(path string) (*rsa.PrivateKey, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %w", err)
	}

	block, _ := pem.Decode(data)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, fmt.Errorf("invalid PEM block or type: %v", block.Type)
	}

	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	return key, nil
}

// EncryptWithPublicKey шифрует данные с помощью открытого ключа.
func EncryptWithPublicKey(data []byte, pub *rsa.PublicKey) ([]byte, error) {
	if pub == nil {
		return nil, errors.New("public key is nil")
	}

	encryptedData, err := rsa.EncryptPKCS1v15(rand.Reader, pub, data)
	if err != nil {
		return nil, fmt.Errorf("encrypt with public key: %w", err)
	}

	return encryptedData, nil
}

// DecryptRSA расшифровывает данные с использованием закрытого ключа (PKCS1v15).
func DecryptRSA(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	if privateKey == nil {
		return nil, errors.New("private key is nil")
	}

	decryptedData, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, data)
	if err != nil {
		return nil, fmt.Errorf("rsa decryption failed: %w", err)
	}

	return decryptedData, nil
}
