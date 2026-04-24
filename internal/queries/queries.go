package queries

import (
	"context"
	"fmt"
	"pwgen/internal/db"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Queries struct {
	db *db.Queries
}

func NewQueries(pool *pgxpool.Pool) *Queries {
	return &Queries{
		db: db.New(pool),
	}
}

func (q *Queries) GetAllVaults(ctx context.Context) ([]db.Vault, error) {
	vaults, err := q.db.GetAllVaults(ctx)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return nil, fmt.Errorf("could not get vaults")
	}
	return vaults, nil
}

func (q *Queries) GetEntryByWebsite(ctx context.Context, website string) (db.VaultEntry, error) {
	entry, err := q.db.GetEntryByWebsite(ctx, website)
	if err != nil {
		fmt.Printf("%s", err.Error())
		return db.VaultEntry{}, nil
	}
	return entry, nil
}

func (q *Queries) InsertVault(ctx context.Context, name, salt string) (db.Vault, error) {
	vault, err := q.db.InsertVault(ctx, db.InsertVaultParams{
		DisplayName: name,
		Salt: salt,
	})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return db.Vault{}, fmt.Errorf("could not save vault")
	}
	return vault, nil
}

func (q *Queries) InsertVaultEntry(ctx context.Context, ciphertext, nonce, website, label string, vaultId pgtype.UUID) (db.VaultEntry, error) {
	entry, err := q.db.InsertVaultEntry(ctx, db.InsertVaultEntryParams{
		Ciphertext: ciphertext,
		Nonce: nonce,
		Website: website,
		Label: label,
		VaultID: vaultId,
	})
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return db.VaultEntry{}, fmt.Errorf("could not inser vault entry")
	}
	return entry, nil
}

func (q *Queries) InsertCurrentVault(ctx context.Context, currentVaultID pgtype.UUID) (db.CurrentVault, error) {
	currentVault, err := q.db.InsertCurrentVault(ctx, currentVaultID)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		return db.CurrentVault{}, fmt.Errorf("could not insert current vault")
	}
	return currentVault, nil
}

func (q *Queries) UpdateCurrentVault(ctx context.Context, currentVaultID pgtype.UUID) (db.CurrentVault, error) {
	currentVault, err := q.db.UpdateCurrentVault(ctx, currentVaultID)
	if err != nil {
		return db.CurrentVault{}, fmt.Errorf("could not update current vault")
	}
	return currentVault, nil
}
