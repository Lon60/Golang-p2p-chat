package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

var publicKey *rsa.PublicKey

func ExportPublicKey() ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("public key not available")
	}

	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return pubBytes, nil
}

func GenerateKeyPairIfNotExists(bits int) error {
	privateKeyPath := "private_key.pem"

	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return GenerateKeyPair(bits)
	}

	_, err := LoadPrivateKey()
	return err
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyPath := "private_key.pem"
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, errors.New("private key file not found")
	}

	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("error decoding private key")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey = &privateKey.PublicKey

	return privateKey, nil
}

func GenerateKeyPair(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	publicKey = &privateKey.PublicKey

	return savePrivateKey("private_key.pem", privateKey)
}

func savePrivateKey(fileName string, key *rsa.PrivateKey) error {
	keyBytes := x509.MarshalPKCS1PrivateKey(key)
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: keyBytes,
	})
	return os.WriteFile(fileName, keyPEM, 0600)
}

func VerifySignature(message string, signature []byte, publicKey *rsa.PublicKey) error {
	hash := sha256.New()
	hash.Write([]byte(message))
	digest := hash.Sum(nil)

	err := rsa.VerifyPSS(publicKey, crypto.SHA256, digest, signature, nil)
	if err != nil {
		return errors.New("signature could not be verified")
	}
	return nil
}

func ImportPublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("error decoding public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("public key is not an RSA key")
	}
}

func DecryptMessage(ciphertext []byte) (string, error) {
	privateKey, err := LoadPrivateKey()
	if err != nil {
		return "", fmt.Errorf("error loading private-key: %v", err)
	}
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("error decrypting message: %v", err)
	}
	return string(plaintext), nil
}
