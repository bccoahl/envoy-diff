# envoy-diff

> CLI tool to diff environment variable sets across `.env` files and remote secrets managers

---

## Installation

```bash
go install github.com/your-org/envoy-diff@latest
```

Or download a pre-built binary from the [Releases](https://github.com/your-org/envoy-diff/releases) page.

---

## Usage

Compare two `.env` files:

```bash
envoy-diff .env.staging .env.production
```

Compare a local `.env` file against a remote secrets manager:

```bash
envoy-diff .env aws://my-app/production
envoy-diff .env vault://secret/my-app
```

**Example output:**

```
~ DATABASE_URL   [changed]
+ NEW_FEATURE_FLAG   (only in production)
- LEGACY_API_KEY     (only in staging)
```

### Supported Sources

| Source | Example |
|--------|---------|
| Local file | `.env`, `.env.production` |
| AWS Secrets Manager | `aws://my-secret-name` |
| HashiCorp Vault | `vault://secret/path` |

### Flags

```
--format    Output format: text, json, yaml (default: text)
--keys-only Show only key names, not values
--no-color  Disable colored output
```

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)