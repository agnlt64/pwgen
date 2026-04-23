package main

import (
	"bufio"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"pwgen/internal/queries"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/argon2"
)

// from Python's string.printable
var PRINTABLE = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ ")
var q *queries.Queries

const (
	MAX_CPU = 10
)

func encryptGCM(plainText, key []byte) (cipherText []byte, nonce []byte, err error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil, err
	}

	nonce = make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, nil, err
	}

	return gcm.Seal(nil, nonce, plainText, nil), nonce, nil
}

func decryptGCM(cipherText, nonce, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherText) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	return gcm.Open(nil, nonce, cipherText, nil)
}

func getMasterPassword() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter master password: ")
	master, err := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return master
}

func randString(size int) string {
	b := make([]rune, size)
	for i := range b {
		// can't crash with rand.Reader according to doc
		idx, _ := rand.Int(rand.Reader, big.NewInt(int64(len(PRINTABLE))))
		b[i] = PRINTABLE[idx.Int64()]
	}
	return string(b)
}

func handleNewVault() {
	ctx := context.Background()
	vaults, err := q.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// todo: this is obviously very dumb
	if len(vaults) >= 1 {
		log.Fatal("a vault already exists")
	}
	salt := randString(16)
	saltB64 := base64.StdEncoding.EncodeToString([]byte(salt))
	fmt.Printf("salt: %s\n", saltB64)
	err = q.InsertVault(ctx, salt)
	if err != nil {
		log.Fatal(err)
	}
}

func handleNewPass() {
	ctx := context.Background()
	vaults, err := q.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}
	vault := vaults[0]
	salt := vault.Salt
	master := getMasterPassword()

	key := argon2.Key([]byte(master), []byte(salt), 3, 32*1024, MAX_CPU, 32)
	// todo: get the size from CLI args
	passwd := randString(10)
	
	cipher, nonce, err := encryptGCM([]byte(passwd), []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new pass: %s\n", passwd)
	entry, err := q.InsertVaultEntry(ctx, base64.StdEncoding.EncodeToString(cipher), base64.StdEncoding.EncodeToString(nonce), "youtube.com", "YT", vault.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("entry cipher: %s\n", entry.Ciphertext)
}

func handleGetPass() {
	ctx := context.Background()
	vaults, err := q.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}
	vault := vaults[0]
	salt := vault.Salt

	entry, err := q.GetEntryByWebsite(ctx, "youtube.com")
	if err != nil {
		log.Fatal(err)
	}
	cipher, err := base64.StdEncoding.DecodeString(entry.Ciphertext)
	if err != nil {
		log.Fatal(err)
	}
	nonce, err := base64.StdEncoding.DecodeString(entry.Nonce)
	if err != nil {
		log.Fatal(err)
	}

	master := getMasterPassword()
	key := argon2.Key([]byte(master), []byte(salt), 3, 32*1024, MAX_CPU, 32)
	plain, err := decryptGCM([]byte(cipher), []byte(nonce), key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decrypted: %s\n", plain)
}

func connect(ctx context.Context, dbURL string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dbURL)

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	return pool
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	if len(os.Args) < 2 {
		log.Fatal("not enough args")
	}
	subCmd := os.Args[1]
	dbURL := os.Getenv("DB_URL")
	pool := connect(context.Background(), dbURL)
	q = queries.NewQueries(pool)

	switch subCmd {
	case "new-vault":
		handleNewVault()
	case "new-pass":
		handleNewPass()
	case "get-pass":
		handleGetPass()
	// todo: help subcommand
	default:
		log.Fatalf("%s is not a valid subcommand", subCmd)
	}
}