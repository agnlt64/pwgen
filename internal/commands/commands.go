package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"pwgen/internal/db"
	"pwgen/internal/security"
	"pwgen/internal/utils"

	"github.com/atotto/clipboard"
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

func checkf(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err.Error())
		os.Exit(1)
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

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	salt := utils.RandString(utils.SALT_LEN)
	displayName := c.args[0]
	vault, err := c.db.InsertVault(ctx, db.InsertVaultParams{
		DisplayName: displayName,
		Salt:        salt,
	})
	checkf(err, "couldn't insert vault")

	_, err = c.db.InsertCurrentVault(ctx, vault.ID)
	checkf(err, "couldn't insert current vault")

	fmt.Printf("Vault %s created successfully! Using it as default vault.\n", displayName)
}

func (c *Commands) UseVault() {
	if len(c.args) != 1 {
		log.Fatal("Error: use-vault command expects a vault name")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := c.args[0]
	vault, err := c.db.GetVaultByName(ctx, name)
	checkf(err, "couldn't get vault by name")

	_, err = c.db.InsertCurrentVault(ctx, vault.ID)
	checkf(err, "couldn't insert current vault")

	fmt.Printf("Using vault %s as default vault.\n", vault.DisplayName)
}

func (c *Commands) ListVaults() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vaults, err := c.db.GetAllVaults(ctx)
	checkf(err, "couldn't get all vaults")

	currentVault, err := c.db.GetCurrentVault(ctx)
	checkf(err, "couldn't insert current vault")

	for idx, vault := range vaults {
		if currentVault.ID == vault.ID {
			fmt.Printf("[%d] %s\n", idx+1, vault.DisplayName)
		} else {
			fmt.Printf(" %d  %s\n", idx+1, vault.DisplayName)
		}
	}
}

func (c *Commands) NewPass() {
	if len(c.args) != 3 {
		// TODO: make length optional
		log.Fatalln("Error: command new-pass expects exactly 3 arguments: new-pass LENGTH WEBSITE_URL WEBSITE_LABEL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vault, err := c.db.GetCurrentVault(ctx)
	checkf(err, "couldn't get current vault")

	salt := vault.Salt
	size, err := strconv.Atoi(c.args[0])
	checkf(err, "invalid size value")

	// TODO: proper URL parsing
	url := c.args[1]
	label := c.args[2]

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)

	passwd := utils.RandString(size)

	cipher, nonce, err := security.Encrypt([]byte(passwd), key)
	checkf(err, "couldn't encrypt password")

	_, err = c.db.InsertVaultEntry(ctx, db.InsertVaultEntryParams{
		Ciphertext: utils.EncodeB64(cipher),
		Nonce:      utils.EncodeB64(nonce),
		Website:    url,
		Label:      label,
		VaultID:    vault.ID,
	})
	checkf(err, "couldn't save password")

	err = clipboard.WriteAll(passwd)
	checkf(err, "couldn't write password to clipboard")

	fmt.Printf("New password for %s was written to clipboard\n", label)
}

func (c *Commands) GetPass() {
	if len(c.args) != 1 {
		log.Fatalln("Error: get-pass command expects exactly 1 argument: WEBSITE_LABEL")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vault, err := c.db.GetCurrentVault(ctx)
	checkf(err, "couldnt' get current vault")

	salt := vault.Salt
	label := c.args[0]

	entry, err := c.db.GetEntryByLabel(ctx, db.GetEntryByLabelParams{
		Label:   label,
		VaultID: vault.ID,
	})
	checkf(err, "invalid label")

	cipher, err := utils.DecodeB64(entry.Ciphertext)
	checkf(err, "couldn't decode base64 literal")

	nonce, err := utils.DecodeB64(entry.Nonce)
	checkf(err, "couldn't decode base64 literal")

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)
	plain, err := security.Decrypt(cipher, nonce, key)
	checkf(err, "invalid password")

	err = clipboard.WriteAll(string(plain))
	checkf(err, "couldn't write password to clipboard")

	fmt.Printf("Password for %s was written to clipboard\n", label)
}
