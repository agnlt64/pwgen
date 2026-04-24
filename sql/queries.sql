-- name: GetAllVaults :many
select *
from vault;

-- name: GetVaultByName :one
select *
from vault
where display_name = $1;

-- name: InsertVault :one
insert into vault (display_name, salt)
values ($1, $2)
returning *;

-- name: GetEntryByWebsite :one
select *
from vault_entry
where website = $1;

-- name: InsertVaultEntry :one
insert into vault_entry (ciphertext, nonce, website, label, vault_id)
values ($1, $2, $3, $4, $5)
returning *;

-- name: GetCurrentVault :one
select *
from current_vault
limit 1;

-- name: InsertCurrentVault :one
insert into current_vault (current_vault_id)
values ($1)
on conflict (singleton) do update set
current_vault_id = $1
returning *;
