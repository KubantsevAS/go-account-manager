package encrypter

import "os"

type Encrypter struct {
	Key string
}

func NewEncrypter() *Encrypter {
	key := os.Getenv("KEY")

	if key == "" {
		panic("parameter KEY is undefined")
	}

	return &Encrypter{
		Key: key,
	}
}

func (enc *Encrypter) Encrypt(plainStr string) string {
	return ""
}

func (enc *Encrypter) Decrypt(plainStr string) string {
	return ""
}
