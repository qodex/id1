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

func encrypt(publicKeyPEM string, data string) ([]byte, error) {
	result := []byte{}

	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return result, fmt.Errorf("invalid key")
	}

	var publicKey *rsa.PublicKey
	if pubKey, err := x509.ParsePKCS1PublicKey(block.Bytes); err == nil {
		publicKey = pubKey
	} else if pubKey, err := x509.ParsePKIXPublicKey(block.Bytes); err == nil {
		publicKey = pubKey.(*rsa.PublicKey)
	} else {
		return result, fmt.Errorf("error parsing key")
	}

	if encrypted, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, []byte(data)); err != nil {
		return []byte{}, err
	} else {
		return encrypted, nil
	}
}

func generateChallenge(id, publicKey string) (string, error) {
	secret := generateSecret(id)
	if encryptedSecret, err := encrypt(publicKey, secret); err != nil {
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
	secret := base64.StdEncoding.EncodeToString(hashBytes)
	return secret
}
