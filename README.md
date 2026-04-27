# pwgen

## How it works
A simple CLI to manage your passwords. Group passwords with vaults. You need at least one vault for the tool to work.

## Build
```console
$ go build .
```

## Usage

### Vaults

Create a vault:
```console
$ ./pwgen vault new VaultName
```

Switch to a vault:
```console
$ ./pwgen vault use VaultName
```

List all vaults:
```console
$ ./pwgen vault list
```

### Passwords

Create a password for a website:
```console
$ ./pwgen pass new https://website.com label
```

Use `--length` to set the password length (default: 25):
```console
$ ./pwgen pass new https://website.com label --length 16
```

Get a password (copied to clipboard):
```console
$ ./pwgen pass get label
```