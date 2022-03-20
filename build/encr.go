package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	b64 "encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func main() {

	data := os.Args[1]

	fmt.Println(completecrypt(data))

}

func completecrypt(data string) string {
	encript := base_ecr(data)
	bytes := []byte(os.Args[2])
	key := hex.EncodeToString(bytes)
	encrypted := encrypt(encript, key)
	encrypted = encrypt(encrypted, key)
	result := fmt.Sprintf("%s", encrypted)
	return result
}

func base_ecr(data string) string {
	sEnc := b64.StdEncoding.EncodeToString([]byte(data))
	uEnc := b64.URLEncoding.EncodeToString([]byte(sEnc))
	return uEnc
}

func encrypt(stringToEncrypt string, keyString string) (encryptedString string) {
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}
