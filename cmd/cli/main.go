package main

import (
	"bufio"
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"flag"
	"io"
	"log"
	"os"
	"strings"
)

func main() {
	var encryptArgs = flag.String("encrypt", "config/app.env.inc", "file path to store")
	var decryptArgs = flag.String("decrypt", "", "encrypted file name with path")
	var key = flag.String("key", "", "encryption key")

	flag.Parse()

	// rx := regexp.MustCompile(`\S`)
	// fileBytes := fileBuffer.Bytes()
	// fileBytes := bytes.TrimSuffix(fileBuffer.Bytes(), []byte(" "))

	*key = strings.TrimSpace(*key)

	if *key == "" {
		log.Fatal("key needs to be present")
	}

	if *decryptArgs != "" {
		plainText, err := decrypt(*decryptArgs, *key)
		if err != nil {
			log.Fatal(err)
		}

		log.Println(plainText)
		return
	}

	if *encryptArgs != "" {
		var fileBuffer = &bytes.Buffer{}
		reader := bufio.NewScanner(os.Stdin)

		for reader.Scan() {
			text := reader.Text()
			fileBuffer.WriteString(text + "\n")

			// truncate white space
			// if ismatch := rx.MatchString(text); !ismatch {
			// 	continue
			// }
			// fmt.Printf("{%s}\n", text)
		}

		err := encrypt(fileBuffer, *encryptArgs, *key)
		if err != nil {
			log.Fatal(err)
		}

		log.Println("encrypted file written to", *encryptArgs)
		return
	}

}

func decrypt(fileName string, key string) (plainText string, err error) {
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

func encrypt(fileBuffer *bytes.Buffer, fileName string, key string) (err error) {
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
	// encrypted = hex.EncodeToString(encryptedBytes)

	err = os.WriteFile(fileName, encryptedBytes, 0644)
	// hashed := sha1.Sum(fileBytes)
	// hashHex := hex.EncodeToString(hashed[:])
	// fmt.Println(hashHex)

	// sha := base64.StdEncoding.EncodeToString([]byte(hashHex))

	// fmt.Println(sha)

	return
}
