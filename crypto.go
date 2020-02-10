package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Encrypt encrypts data using RB's API method with aes sha256
func Encrypt(data string, key string) string {
	codeData := data + "\x00\x10"

	m1 := sha256.New()
	m1.Write([]byte(key))

	b1 := m1.Sum(nil)

	m2 := sha256.New()
	m2.Write([]byte(data))

	b2 := m2.Sum(nil)

	block, err := aes.NewCipher([]byte(b1))

	if err != nil {
		log.Error(err)
	}

	if len(codeData)%aes.BlockSize != 0 {
		dataLen := len(codeData)

		remaining := dataLen % aes.BlockSize

		padLength := aes.BlockSize - remaining

		codeData = codeData + strings.Repeat("\x10", padLength)
	}

	codeByteData := []byte(codeData)

	ciphertext := make([]byte, aes.BlockSize+len(codeByteData))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		log.Error(err)
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], codeByteData)

	return string(b2) + string(ciphertext)
}

// Decrypt decrypts data using RB's API method with aes sha256
func Decrypt(encryptedData string, key string) string {
	iv := []byte(encryptedData[32:48])

	encryptedByteData := []byte(encryptedData[48:])

	m := sha256.New()
	m.Write([]byte(key))

	symmetricKey := m.Sum(nil)[:32]
	block, err := aes.NewCipher(symmetricKey)

	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptedByteData, encryptedByteData)

	return string(encryptedByteData)
}
