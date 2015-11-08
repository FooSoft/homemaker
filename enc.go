package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/andrewrynhard/go-mask"
)

func decrypt(cipherstring string, keystring string) string {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	// Get the 16 byte IV
	iv := ciphertext[:aes.BlockSize]

	// Remove the IV from the ciphertext
	ciphertext = ciphertext[aes.BlockSize:]

	// Return a decrypted stream
	stream := cipher.NewCFBDecrypter(block, iv)

	// Decrypt bytes from ciphertext
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(plainstring, keystring string) string {
	// Byte array of the string
	plaintext := []byte(plainstring)

	// Key
	key := []byte(keystring)

	// Create the AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	// Empty array of 16 + plaintext length
	// Include the IV at the beginning
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// Slice of first 16 bytes
	iv := ciphertext[:aes.BlockSize]

	// Write 16 rand bytes to fill iv
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	// Return an encrypted stream
	stream := cipher.NewCFBEncrypter(block, iv)

	// Encrypt bytes from plaintext to ciphertext
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return string(ciphertext)
}

func writeToFile(data, file string) {
	ioutil.WriteFile(file, []byte(data), 0666)
}

func readFromFile(file string) ([]byte, error) {
	data, err := ioutil.ReadFile(file)
	return data, err
}

func isFile(file string) (bool, error) {
	s, err := os.Stat(file)

	if os.IsNotExist(err) {
		return false, err
	}

	if s.IsDir() {
		return false, err
	}

	return true, nil
}

func readKey() (string, error) {
	maskedReader := mask.NewMaskedReader()

	key, err := maskedReader.GetInputConfirmMasked()
	if err != nil {
		return "", err
	}

	hasher := md5.New()
	hasher.Write([]byte(key))

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func encryptFile(file string, key string) error {
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	encrypted := encrypt(string(content), string(key))
	writeToFile(encrypted, file+".enc")

	if prompt(strings.Join([]string{"Remove", file, "?"}, " ")) {
		os.Remove(file)
	}

	return nil
}

func decryptFile(file string, key string) error {
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	decrypted := decrypt(string(content), string(key))
	writeToFile(decrypted, file[:len(file)-4])

	return nil
}

func encryptEnc(param string, conf *config, key string) error {
	sourceFile := filepath.Join(conf.srcDir, param)

	exists, err := isFile(sourceFile)
	if !exists {
		return err
	}

	encryptFile(sourceFile, key)

	return nil
}

func processEnc(param string, conf *config, key string) error {
	sourceFile := filepath.Join(conf.srcDir, param) + ".enc"

	exists, err := isFile(sourceFile)
	if !exists {
		return err
	}

	decryptFile(sourceFile, key)

	return nil
}
