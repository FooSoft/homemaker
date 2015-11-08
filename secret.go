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
	ciphertext := []byte(cipherstring)

	key := []byte(keystring)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	iv := ciphertext[:aes.BlockSize]

	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext)
}

func encrypt(plainstring, keystring string) string {
	plaintext := []byte(plainstring)

	key := []byte(keystring)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	iv := ciphertext[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)

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

func getMaskedInput() ([]byte, error) {
	maskedReader := mask.NewMaskedReader()

	key, err := maskedReader.GetInputConfirmMasked()
	if err != nil {
		return nil, err
	}

	return key, nil
}

func keyFromPasssword(password []byte) string {
	hasher := md5.New()

	hasher.Write([]byte(password))

	return hex.EncodeToString(hasher.Sum(nil))
}

func encryptFile(file string, key string, conf *config) error {
	content, err := readFromFile(file)
	if err != nil {
		return err
	}

	encrypted := encrypt(string(content), string(key))
	writeToFile(encrypted, file+".enc")

	if conf.remove || prompt(strings.Join([]string{"Remove", file, "?"}, " ")) {
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

func encryptSecret(param string, conf *config, key string) error {
	sourceFile := filepath.Join(conf.srcDir, param)

	exists, err := isFile(sourceFile)
	if !exists {
		return err
	}

	encryptFile(sourceFile, key, conf)

	return nil
}

func decryptSecret(param string, conf *config, key string) error {
	sourceFile := filepath.Join(conf.srcDir, param) + ".enc"

	exists, err := isFile(sourceFile)
	if !exists {
		return err
	}

	decryptFile(sourceFile, key)

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
