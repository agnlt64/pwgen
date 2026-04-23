-- name: GetAllVaults :many
select *
from vault;

-- name: GetEntryByWebsite :one
select *
from vault_entry
where website = $1;

-- name: InsertVault :one
insert into vault (salt)
values ($1)
returning *;

-- name: InsertVaultEntry :one
insert into vault_entry (ciphertext, nonce, website, label, vault_id)
values ($1, $2, $3, $4, $5)
returning *;