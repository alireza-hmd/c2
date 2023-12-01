package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"io/ioutil"
	"log"
)

func GenerateToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}

func ReadFile(path string) ([]byte, error) {
	text, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error reading file")
	}
	return text, nil
}

func Encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error generating cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error generating gcm")
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		log.Println(err)
		return nil, errors.New("nonce error")
	}

	cipherText := gcm.Seal(nonce, nonce, data, nil)
	return cipherText, nil
}

func Decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error generating cipher")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error generating gcm")
	}

	nonce := data[:gcm.NonceSize()]
	cipherText := data[gcm.NonceSize():]
	text, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		log.Println(err)
		return nil, errors.New("error decrypting data")
	}
	return text, nil
}

func StrToByte(str string) []byte {
	byteString := make([]byte, len(str))
	copy(byteString, str)
	return byteString
}
