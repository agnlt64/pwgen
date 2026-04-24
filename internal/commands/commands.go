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
	args 	[]string
}

func NewCommands(queries *queries.Queries, args []string) *Commands {
	return &Commands{
		queries: queries,
		args: args,
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func Usage() {
	// TODO: move to actual subcommands, e.g
	// vault [new|list|use]
	// vault new [name], vault use [name], vault list
	fmt.Println("Usage:")
	fmt.Println("    new-vault [NAME] - Create a new vault")
	fmt.Println("    use-vault [NAME] - Use a specific vault")
	fmt.Println("    list-vaults      - List all vaults")
	fmt.Println("")
	fmt.Println("    new-pass [WEBSITE] - Create a new password (TODO: support website)")
	fmt.Println("    get-pass [WEBSITE] - Get the password for given WEBSITE (TODO: support website)")
	fmt.Println("")
	fmt.Println("    help - Print this help message")
}

func (c *Commands) NewVault() {
	if len(c.args) != 1 {
		log.Fatal("Error: new-vault command expects a vault name")
	}

	ctx := context.Background()
	salt := utils.RandString(utils.SALT_LEN)
	displayName := c.args[0]
	vault, err := c.queries.InsertVault(ctx, displayName, salt)
	check(err)

	_, err = c.queries.UpdateCurrentVault(ctx, vault.ID)
	if err != nil {
		// if we're here, current_vault table is empty, so create it
		_, err = c.queries.InsertCurrentVault(ctx, vault.ID)
		check(err)
	}
	fmt.Printf("Vault %s created successfully! Using it as default vault.\n", displayName)
}

func (c *Commands) UseVault() {
	if len(c.args) != 1 {
		log.Fatal("Error: use-vault command expects a vault name")
	}

	ctx := context.Background()
	name := c.args[0]
	vault, err := c.queries.GetVaultByName(ctx, name)
	check(err)

	_, err = c.queries.UpdateCurrentVault(ctx, vault.ID)
	check(err)
	fmt.Printf("Using vault %s as default vault.\n", vault.DisplayName)
}

func (c *Commands) ListVaults() {
	ctx := context.Background()
	vaults, err := c.queries.GetAllVaults(ctx)
	check(err)

	currentVault, err := c.queries.GetCurrentVault(ctx)
	check(err)

	for idx, vault := range vaults {
		if currentVault.CurrentVaultID == vault.ID {
			fmt.Printf("[%d] %s\n", idx, vault.DisplayName)
		} else {
			fmt.Printf(" %d  %s\n", idx, vault.DisplayName)
		}
	}
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
