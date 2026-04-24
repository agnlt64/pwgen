package utils

import (
	"bufio"
	"log"
	"fmt"
	"os"
	"math/big"
	"crypto/rand"
)

// from Python's string.printable
var PRINTABLE = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")

const (
	MAX_CPU = 8
)

func GetMasterPassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter master password: ")
	master, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return master
}

func RandString(size int) string {
	b := make([]rune, size)
	for i := range b {
		// can't crash with rand.Reader according to doc
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(PRINTABLE))))
		b[i] = PRINTABLE[idx.Int64()]
	}
	return string(b)
}
