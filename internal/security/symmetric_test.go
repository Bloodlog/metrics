package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEncryptDecryptRSA(t *testing.T) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	publicKey := &privateKey.PublicKey

	data := []byte("hello world")

	encrypted, err := EncryptWithPublicKey(data, publicKey)
	require.NoError(t, err)

	decrypted, err := DecryptRSA(privateKey, encrypted)
	require.NoError(t, err)

	require.Equal(t, data, decrypted)
}

func TestLoadAndUseRSAKeysFromFiles(t *testing.T) {
	tmpDir := t.TempDir()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privKeyPath := filepath.Join(tmpDir, "key.pem")
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPem := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes})
	require.NoError(t, os.WriteFile(privKeyPath, privPem, 0o600))

	template := &x509.Certificate{SerialNumber: bigInt(1)}
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	require.NoError(t, err)
	certPem := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certBytes})
	certPath := filepath.Join(tmpDir, "cert.pem")
	require.NoError(t, os.WriteFile(certPath, certPem, 0o600))

	loadedPriv, err := LoadRSAPrivateKeyFromFile(privKeyPath)
	require.NoError(t, err)
	require.Equal(t, privateKey.D, loadedPriv.D)

	loadedPub, err := LoadRSAPublicKeyFromCert(certPath)
	require.NoError(t, err)
	require.Equal(t, privateKey.PublicKey.E, loadedPub.E)
	require.Equal(t, privateKey.PublicKey.N.Cmp(loadedPub.N), 0)

	msg := []byte("test message")
	encrypted, err := EncryptWithPublicKey(msg, loadedPub)
	require.NoError(t, err)

	decrypted, err := DecryptRSA(loadedPriv, encrypted)
	require.NoError(t, err)
	require.Equal(t, msg, decrypted)
}

func bigInt(val int64) *big.Int {
	return new(big.Int).SetInt64(val)
}
