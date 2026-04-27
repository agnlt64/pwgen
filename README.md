# pwgen

## How it works
A simple CLI to manage your passwords. Group passwords with vaults. You need at least one vault for the tool to work.

## Build
```console
$ go build .
```

## Usage
Create a vault:
```console
$ ./pwgen new-vault VaultName
```

List all vaults:
```console
$ ./pwgen list-vaults
```

Create a 10-chars password for a specific website:
```console
$ ./pwgen new-pass 10 https://website.com website
```

Get a password for a given website (website label, not the URL):
```console
$ ./pwgen get-pass website
```