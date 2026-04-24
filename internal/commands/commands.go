package commands

import (
	"context"
	"encoding/base64"
	"log"
	"fmt"

	"pwgen/internal/security"
	"pwgen/internal/queries"
	"pwgen/internal/utils"
	"golang.org/x/crypto/argon2"
)

type Commands struct {
	queries *queries.Queries
}

func NewCommands(queries *queries.Queries) *Commands {
	return &Commands{
		queries: queries,
	}
}

func (c *Commands) NewVault() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	// todo: this is obviously very dumb
	if len(vaults) >= 1 {
		log.Fatal("a vault already exists")
	}
	salt := utils.RandString(16)
	saltB64 := base64.StdEncoding.EncodeToString([]byte(salt))
	fmt.Printf("salt: %s\n", saltB64)
	err = c.queries.InsertVault(ctx, salt)
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Commands) NewPass() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}
	vault := vaults[0]
	salt := vault.Salt
	master := utils.GetMasterPassword()

	key := argon2.Key([]byte(master), []byte(salt), 3, 32*1024, utils.MAX_CPU, 32)
	// todo: get the size from CLI args
	passwd := utils.RandString(10)
	
	cipher, nonce, err := security.Encrypt([]byte(passwd), []byte(key))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new pass: %s\n", passwd)
	entry, err := c.queries.InsertVaultEntry(ctx, base64.StdEncoding.EncodeToString(cipher), base64.StdEncoding.EncodeToString(nonce), "youtube.com", "YT", vault.ID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("entry cipher: %s\n", entry.Ciphertext)
}

func (c *Commands) GetPass() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	if err != nil {
		log.Fatal(err)
	}
	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}
	vault := vaults[0]
	salt := vault.Salt

	entry, err := c.queries.GetEntryByWebsite(ctx, "youtube.com")
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

	master := utils.GetMasterPassword()
	key := argon2.Key([]byte(master), []byte(salt), 3, 32*1024, utils.MAX_CPU, 32)
	plain, err := security.Decrypt([]byte(cipher), []byte(nonce), key)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("decrypted: %s\n", plain)
}
