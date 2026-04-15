# vaultpipe

> A CLI tool for injecting secrets from HashiCorp Vault into local dev environments via `.env` files.

---

## Installation

```bash
go install github.com/yourusername/vaultpipe@latest
```

Or download a pre-built binary from the [releases page](https://github.com/yourusername/vaultpipe/releases).

---

## Usage

Point `vaultpipe` at a Vault path and it will write the secrets to a `.env` file in your project directory.

```bash
# Authenticate with Vault first (uses VAULT_ADDR and VAULT_TOKEN env vars)
export VAULT_ADDR="https://vault.example.com"
export VAULT_TOKEN="s.xxxxxxxxxxxxxxxx"

# Pull secrets from a Vault KV path and write to .env
vaultpipe inject --path secret/data/myapp/dev --output .env

# Preview secrets without writing to disk
vaultpipe inject --path secret/data/myapp/dev --dry-run

# Append secrets to an existing .env file
vaultpipe inject --path secret/data/myapp/dev --output .env --append
```

**Example `.env` output:**

```
DB_HOST=localhost
DB_PASSWORD=supersecret
API_KEY=abc123
```

### Flags

| Flag | Description | Default |
|------|-------------|----------|
| `--path` | Vault KV secret path | *(required)* |
| `--output` | Output file path | `.env` |
| `--append` | Append to existing file instead of overwriting | `false` |
| `--dry-run` | Print secrets to stdout without writing to disk | `false` |

> **Note:** `--dry-run` and `--append` are mutually exclusive. If both are set, `--dry-run` takes precedence.

---

## Requirements

- Go 1.21+
- A running [HashiCorp Vault](https://www.vaultproject.io/) instance
- A valid `VAULT_ADDR` and `VAULT_TOKEN` (or other supported auth method)

---

## License

[MIT](LICENSE) © 2024 yourusername
