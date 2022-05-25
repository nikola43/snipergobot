package encryptionutils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
)

func rKey(filename string) ([]byte, error) {
	key, err := ioutil.ReadFile(filename)
	if err != nil {
		return key, err
	}
	block, _ := pem.Decode(key)
	return block.Bytes, nil
}

func cKey() []byte {
	genkey := make([]byte, 16)
	_, err := rand.Read(genkey)
	if err != nil {
		log.Fatalf("failed to read key: %s", err)
	}
	return genkey
}

func sKey(filename string, key []byte) {
	block := &pem.Block{
		Type:  "AES KEY",
		Bytes: key,
	}
	err := ioutil.WriteFile(filename, pem.EncodeToMemory(block), 9854)
	if err != nil {
		log.Fatalf("Failed tio save the key %s: %s", filename, err)
	}
}

func aesKey(keyFile string) []byte {
	file := fmt.Sprintf(keyFile)
	key, err := rKey(file)
	if err != nil {
		log.Println("Create a new AES KEY")
		key = cKey()
		sKey(file, key)
	}
	return key
}

func createCipher(keyFile string) cipher.Block {
	c, err := aes.NewCipher(aesKey(keyFile))
	if err != nil {
		log.Fatalf("failed to create aes  %s", err)
	}
	return c
}

func encryption(plainText string, filename string, keyFile string, encryptionKey []byte) {
	bytes := []byte(plainText)
	blockCipher := createCipher(keyFile)
	stream := cipher.NewCTR(blockCipher, encryptionKey)
	stream.XORKeyStream(bytes, bytes)
	err := ioutil.WriteFile(fmt.Sprintf(filename), bytes, 0644)
	if err != nil {
		log.Fatalf("writing encryption file %s", err)
	}
}
func decryption(filename string, keyFile string, encryptionKey []byte) []byte {
	bytes, err := ioutil.ReadFile(fmt.Sprintf(filename))
	if err != nil {
		log.Fatalf("Reading encrypted file %s", err)
	}
	blockCipher := createCipher(keyFile)
	stream := cipher.NewCTR(blockCipher, encryptionKey)
	stream.XORKeyStream(bytes, bytes)
	return bytes
}
