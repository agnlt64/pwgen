package commands

import (
	"context"
	"log"
	"fmt"

	"pwgen/internal/security"
	"pwgen/internal/queries"
	"pwgen/internal/utils"
)

type Commands struct {
	queries *queries.Queries
}

func NewCommands(queries *queries.Queries) *Commands {
	return &Commands{
		queries: queries,
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c *Commands) NewVault() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	check(err)

	// todo: this is obviously very dumb
	if len(vaults) >= 1 {
		log.Fatal("a vault already exists")
	}

	salt := utils.RandString(16)
	err = c.queries.InsertVault(ctx, salt)
	check(err)
}

func (c *Commands) NewPass() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	check(err)

	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}

	vault := vaults[0]
	salt := vault.Salt
	master := utils.GetMasterPassword()

	key := utils.Argon2id(master, salt)
	// todo: get the size from CLI args
	passwd := utils.RandString(10)

	cipher, nonce, err := security.Encrypt([]byte(passwd), key)
	check(err)

	fmt.Printf("new pass: %s\n", passwd)
	entry, err := c.queries.InsertVaultEntry(ctx, utils.EncodeB64(cipher), utils.EncodeB64(nonce), "youtube.com", "YT", vault.ID)
	check(err)

	fmt.Printf("entry cipher: %s\n", entry.Ciphertext)
}

func (c *Commands) GetPass() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	check(err)

	if len(vaults) > 1 {
		log.Fatal("more than one vault is not allowed yet")
	}

	vault := vaults[0]
	salt := vault.Salt

	entry, err := c.queries.GetEntryByWebsite(ctx, "youtube.com")
	check(err)

	cipher := utils.DecodeB64(entry.Ciphertext)
	nonce := utils.DecodeB64(entry.Nonce)

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)
	plain, err := security.Decrypt([]byte(cipher), []byte(nonce), key)
	check(err)

	fmt.Printf("decrypted: %s\n", plain)
}
