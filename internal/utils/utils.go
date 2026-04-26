package utils

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"math/big"
	"os"

	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

// from Python's string.printable
var PRINTABLE = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")

const (
	MAX_CPU  = 8
	SALT_LEN = 16
)

func GetMasterPassword() string {
	fmt.Print("Enter master password: ")
	bytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println()
	return string(bytes)
}

func RandString(size int) string {
	b := make([]rune, size)
	for i := range b {
		idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(PRINTABLE))))
		if err != nil {
			// can't crash with rand.Reader according to doc
			panic("unreachable")
		}
		b[i] = PRINTABLE[idx.Int64()]
	}
	return string(b)
}

func Argon2id(password, salt string) []byte {
	key := argon2.Key([]byte(password), []byte(salt), 3, 32*1024, MAX_CPU, 32)
	return key
}

func EncodeB64(text []byte) string {
	return base64.StdEncoding.EncodeToString(text)
}

func DecodeB64(text string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(text)
}
