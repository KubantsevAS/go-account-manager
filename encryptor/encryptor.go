package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"
)

type Encryptor struct {
	Key string
}

func NewEncryptor() *Encryptor {
	key := os.Getenv("KEY")

	if key == "" {
		panic("parameter KEY is undefined")
	}

	return &Encryptor{
		Key: key,
	}
}

func (enc *Encryptor) Encrypt(plainStr []byte) []byte {
	// * Standard authenticated encryption approach using AES-GCM.
	// Creates an AES cipher block from the encryption key

	block, err := aes.NewCipher([]byte(enc.Key))
	panicOnError(err)

	// Wraps the block cipher in GCM mode (Galois Counter Mode)
	// GCM provides authenticated encryption (confidentiality and authenticity)
	// Produces a cipher object for encryption/decryption

	aesGSM, err := cipher.NewGCM(block)
	panicOnError(err)

	// Allocates a nonce (number used once) of the required size
	// Fills it with cryptographically secure random bytes
	// The nonce must be unique per encryption with the same key

	nonce := make([]byte, aesGSM.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	panicOnError(err)

	return aesGSM.Seal(nonce, nonce, plainStr, nil)
}

func (enc *Encryptor) Decrypt(encryptedStr []byte) []byte {
	block, err := aes.NewCipher([]byte(enc.Key))
	panicOnError(err)

	aesGSM, err := cipher.NewGCM(block)
	panicOnError(err)

	nonceSize := aesGSM.NonceSize()

	nonce, cipherText := encryptedStr[:nonceSize], encryptedStr[nonceSize:]

	plainText, err := aesGSM.Open(nil, nonce, cipherText, nil)
	panicOnError(err)

	return plainText
}

func panicOnError(err error) {
	if err != nil {
		panic(err.Error())
	}
}
