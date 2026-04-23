package main

import (
	"fmt"
	"math/rand"
	"os"
	"log"
	"strconv"
)

// from Python's string.printable
var PRINTABLE = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")

func genPassword(size int) string {
	b := make([]rune, size)
    for i := range b {
        b[i] = PRINTABLE[rand.Intn(len(PRINTABLE))]
    }
    return string(b)
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}

	sizeStr := os.Args[1]
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		log.Fatalf("%s is not a valid integer", sizeStr)
	}
	
	if size < 0 {
		log.Fatal("size must be greater than 0")
	}
	passwd := genPassword(size)
	fmt.Println(passwd)
}