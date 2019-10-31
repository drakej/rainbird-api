package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
)

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

		remaining := aes.BlockSize - dataLen
		padLength := remaining % aes.BlockSize

		codeData = codeData + strings.Repeat("\x10", padLength)
	}

	log.Error(hex.EncodeToString([]byte(codeData)))

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

func Decrypt(encryptedData string, key string) string {

	log.Info(hex.EncodeToString([]byte(encryptedData)))
	iv := []byte(encryptedData[32:48])

	encryptedByteData := []byte(encryptedData[48:])

	m := sha256.New()
	m.Write([]byte(key))

	symmetricKey := m.Sum(nil)[:32]

	log.Info(hex.EncodeToString(iv))
	log.Info(hex.EncodeToString(encryptedByteData))

	block, err := aes.NewCipher(symmetricKey)

	if err != nil {
		panic(err)
	}

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(encryptedByteData, encryptedByteData)

	return string(encryptedByteData)
}
