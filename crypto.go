package id1

import (
	"time"

	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func encryptWithPubKey(publicKeyPEM string, message string) ([]byte, error) {
	if block, _ := pem.Decode([]byte(publicKeyPEM)); block == nil {
		return []byte{}, fmt.Errorf("invalid key")
	} else if publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes); err != nil {
		return []byte{}, err
	} else if encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(message)); err != nil {
		return []byte{}, err
	} else {
		return encrypted, nil
	}
}

func generateChallenge(id, publicKey string) (string, error) {
	secret := generateSecret(id)
	if encryptedSecret, err := encryptWithPubKey(publicKey, secret); err != nil {
		return "", err
	} else {
		challenge := base64.StdEncoding.EncodeToString(encryptedSecret)
		return challenge, nil
	}
}

func generateSecret(id string) string {
	str := fmt.Sprintf("for user %s, this string changes once a day: %d", id, time.Now().Day())
	hash := sha256.New()
	hash.Write([]byte(str))
	hashBytes := hash.Sum(nil)
	secret := string(hashBytes)
	return secret
}
