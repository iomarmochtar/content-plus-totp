package enc

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

const keyMinSize = 32

func getGcm(key string) (cipher.AEAD, error) {
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	return gcm, nil
}

// adjustKeySize adjusting key size with min character using right padding by space character
func adjustKeySize(key string, size int) string {
	if len(key) > size {
		return key
	}

	return fmt.Sprintf(fmt.Sprintf("%s%%%ds", key, size-len(key)), "")
}

func Decrypt(key, content string) (string, error) {
	raw, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	key = adjustKeySize(key, keyMinSize)

	gcm, err := getGcm(key)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if nonceSize > len(raw) {
		return "", fmt.Errorf(
			"not a valid length decoded b64 value (%d) with nonce (%d), possibly not a valid encrypted content",
			len(raw), nonceSize)
	}

	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	nonce, raw = raw[:nonceSize], raw[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, raw, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

func Encrypt(key, content string) (string, error) {
	key = adjustKeySize(key, keyMinSize)
	c, err := aes.NewCipher([]byte(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce := make([]byte, nonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	result := gcm.Seal(nonce, nonce, []byte(content), nil)

	return base64.StdEncoding.EncodeToString(result), nil
}
