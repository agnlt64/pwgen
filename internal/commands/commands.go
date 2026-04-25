package commands

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"pwgen/internal/db"
	"pwgen/internal/security"
	"pwgen/internal/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Commands struct {
	args []string
	db   *db.Queries
}

func NewCommands(pool *pgxpool.Pool, args []string) *Commands {
	return &Commands{
		args: args,
		db:   db.New(pool),
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
	fmt.Println("    new-pass [SIZE] [URL] [LABEL] - Create a new password")
	fmt.Println("    get-pass [LABEL] 			   - Get the password for website associated with LABEL")
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
	vault, err := c.db.InsertVault(ctx, db.InsertVaultParams{
		DisplayName: displayName,
		Salt:        salt,
	})
	check(err)

	_, err = c.db.InsertCurrentVault(ctx, vault.ID)
	check(err)
	fmt.Printf("Vault %s created successfully! Using it as default vault.\n", displayName)
}

func (c *Commands) UseVault() {
	if len(c.args) != 1 {
		log.Fatal("Error: use-vault command expects a vault name")
	}

	ctx := context.Background()
	name := c.args[0]
	vault, err := c.db.GetVaultByName(ctx, name)
	check(err)

	_, err = c.db.InsertCurrentVault(ctx, vault.ID)
	check(err)
	fmt.Printf("Using vault %s as default vault.\n", vault.DisplayName)
}

func (c *Commands) ListVaults() {
	ctx := context.Background()
	vaults, err := c.db.GetAllVaults(ctx)
	check(err)

	currentVault, err := c.db.GetCurrentVault(ctx)
	check(err)

	for idx, vault := range vaults {
		if currentVault.CurrentVaultID == vault.ID {
			fmt.Printf("[%d] %s\n", idx+1, vault.DisplayName)
		} else {
			fmt.Printf(" %d  %s\n", idx+1, vault.DisplayName)
		}
	}
}

func (c *Commands) NewPass() {
	if len(c.args) != 3 {
		// todo: make length optional
		log.Fatalln("Error: command new-pass expects exactly 3 arguments: new-pass LENGTH WEBSITE_URL WEBSITE_LABEL")
	}

	ctx := context.Background()
	currentVault, err := c.db.GetCurrentVault(ctx)
	check(err)

	vault, err := c.db.GetVaultById(ctx, currentVault.CurrentVaultID)
	check(err)

	salt := vault.Salt
	size, err := strconv.Atoi(c.args[0])
	check(err)

	// todo: proper URL parsing
	url := c.args[1]
	label := c.args[2]

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)

	passwd := utils.RandString(size)

	cipher, nonce, err := security.Encrypt([]byte(passwd), key)
	check(err)

	_, err = c.db.InsertVaultEntry(ctx, db.InsertVaultEntryParams{
		Ciphertext: utils.EncodeB64(cipher),
		Nonce:      utils.EncodeB64(nonce),
		Website:    url,
		Label:      label,
		VaultID:    vault.ID,
	})
	check(err)

	fmt.Printf("new pass for %s: %s\n", url, passwd)
}

func (c *Commands) GetPass() {
	if len(c.args) != 1 {
		log.Fatalln("Error: get-pass command expects exactly 1 argument: WEBSITE_LABEL")
	}

	ctx := context.Background()
	currentVault, err := c.db.GetCurrentVault(ctx)
	check(err)

	vault, err := c.db.GetVaultById(ctx, currentVault.CurrentVaultID)
	check(err)

	salt := vault.Salt
	label := c.args[0]

	entry, err := c.db.GetEntryByLabel(ctx, db.GetEntryByLabelParams{
		Label:   label,
		VaultID: vault.ID,
	})
	check(err)

	cipher := utils.DecodeB64(entry.Ciphertext)
	nonce := utils.DecodeB64(entry.Nonce)

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)
	plain, err := security.Decrypt(cipher, nonce, key)
	check(err)

	fmt.Printf("decrypted: %s\n", plain)
}
