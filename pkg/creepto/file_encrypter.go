package creepto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"log"
	"os"
)

func Decrypt(fileName string, key string) (plainText string, err error) {
	fileBytes, err := os.ReadFile(fileName)
	if err != nil {
		return plainText, err
	}

	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	var gcm cipher.AEAD
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := fileBytes[:gcm.NonceSize()]
	fileBytes = fileBytes[gcm.NonceSize():]

	var plainTextBytes []byte

	plainTextBytes, err = gcm.Open(nil, nonce, fileBytes, nil)
	if err != nil {
		return
	}

	plainText = string(plainTextBytes)
	return
}

func Encrypt(fileBuffer *bytes.Buffer, fileName string, key string) (err error) {
	fileBytes := bytes.TrimSuffix(fileBuffer.Bytes(), []byte(" "))

	var block cipher.Block
	block, err = aes.NewCipher([]byte(key))
	if err != nil {
		return
	}

	var gcm cipher.AEAD // AD stands for associate data
	gcm, err = cipher.NewGCM(block)
	if err != nil {
		return
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Fatalf("nonce  err: %v", err.Error())
	}

	encryptedBytes := gcm.Seal(nonce, nonce, fileBytes, nil)
	err = os.WriteFile(fileName, encryptedBytes, 0644)

	return
}
