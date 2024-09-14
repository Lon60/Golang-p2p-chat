package security

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
)

var publicKey *rsa.PublicKey

func ExportPublicKey() ([]byte, error) {
	if publicKey == nil {
		return nil, errors.New("Public Key nicht verfügbar")
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

// GenerateKeyPairIfNotExists prüft, ob das Schlüsselpaar existiert, und generiert es falls nicht.
func GenerateKeyPairIfNotExists(bits int) error {
	privateKeyPath := "private_key.pem"

	// Prüfen, ob der private Schlüssel bereits existiert
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		// Falls nicht vorhanden, generiere ein neues Schlüsselpaar
		return GenerateKeyPair(bits)
	}

	// Lade den vorhandenen Private Key, falls er existiert
	_, err := LoadPrivateKey()
	return err
}

func LoadPrivateKey() (*rsa.PrivateKey, error) {
	// Überprüfe, ob die Datei existiert
	privateKeyPath := "private_key.pem"
	if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
		return nil, errors.New("Private Key Datei nicht gefunden")
	}

	// Lade die Private Key Datei
	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	// PEM-Dekodierung
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("Fehler beim Dekodieren des Private Keys")
	}

	// Parse den Private Key
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return privateKey, nil
}

func GenerateKeyPair(bits int) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	publicKey = &privateKey.PublicKey

	// Speichere den Private Key
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

	// Verifiziere die Signatur
	err := rsa.VerifyPSS(publicKey, crypto.SHA256, digest, signature, nil)
	if err != nil {
		return errors.New("Signatur konnte nicht verifiziert werden")
	}
	return nil
}

func ImportPublicKey(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("Fehler beim Dekodieren des Public Keys")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("Public Key ist kein RSA-Schlüssel")
	}
}
