package commands

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"pwgen/internal/db"
	"pwgen/internal/security"
	"pwgen/internal/utils"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
)

func checkf(err error, msg string) {
	if err != nil {
		fmt.Printf("%s: %v\n", msg, err.Error())
		os.Exit(1)
	}
}

func NewVault(queries *db.Queries, cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	salt := utils.RandString(utils.SALT_LEN)
	displayName := args[0]
	vault, err := queries.InsertVault(ctx, db.InsertVaultParams{
		DisplayName: displayName,
		Salt:        salt,
	})
	checkf(err, "couldn't insert vault")

	_, err = queries.InsertCurrentVault(ctx, vault.ID)
	checkf(err, "couldn't insert current vault")

	fmt.Printf("Vault %s created successfully! Using it as default vault.\n", displayName)
}

func UseVault(queries *db.Queries, cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		log.Fatal("Error: use-vault command expects a vault name")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	name := args[0]
	vault, err := queries.GetVaultByName(ctx, name)
	checkf(err, "couldn't get vault by name")

	_, err = queries.InsertCurrentVault(ctx, vault.ID)
	checkf(err, "couldn't insert current vault")

	fmt.Printf("Using vault %s as default vault.\n", vault.DisplayName)
}

func ListVaults(queries *db.Queries, cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vaults, err := queries.GetAllVaults(ctx)
	checkf(err, "couldn't get all vaults")

	currentVault, err := queries.GetCurrentVault(ctx)
	checkf(err, "couldn't insert current vault")

	for idx, vault := range vaults {
		if currentVault.ID == vault.ID {
			fmt.Printf("[%d] %s\n", idx+1, vault.DisplayName)
		} else {
			fmt.Printf(" %d  %s\n", idx+1, vault.DisplayName)
		}
	}
}

func NewPass(queries *db.Queries, cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vault, err := queries.GetCurrentVault(ctx)
	checkf(err, "couldn't get current vault")

	salt := vault.Salt
	length, _ := cmd.Flags().GetInt("length")

	// TODO: proper URL parsing
	url := args[0]
	label := args[1]

	master := utils.GetMasterPassword()
	key := utils.Argon2id(master, salt)

	passwd := utils.RandString(length)

	cipher, nonce, err := security.Encrypt([]byte(passwd), key)
	checkf(err, "couldn't encrypt password")

	_, err = queries.InsertVaultEntry(ctx, db.InsertVaultEntryParams{
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

func GetPass(queries *db.Queries, cmd *cobra.Command, args []string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	vault, err := queries.GetCurrentVault(ctx)
	checkf(err, "couldnt' get current vault")

	salt := vault.Salt
	label := args[0]

	entry, err := queries.GetEntryByLabel(ctx, db.GetEntryByLabelParams{
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
