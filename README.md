# curlsc — Concurrent URL Status Checker

A fast, concurrent CLI utility for real-time endpoint monitoring and network diagnostics.

## Features

- **Concurrent Execution**: Uses Go's goroutines to check multiple URLs simultaneously.
- **Worker Pool**: Configurable concurrency levels to prevent resource exhaustion.
- **Multiple Inputs**: Supports command-line arguments, file input (`--file`), and stdin piping.
- **Graceful Error Handling**: Classifies errors into `TIMEOUT`, `DNS_ERROR`, `CONNECTION_REFUSED`, etc.
- **Rich Output**: Choose between Table (default), JSON, or CSV formats.
- **Color Coding**: Instant visual feedback in the terminal.
- **CI/CD Friendly**: Returns exit code 1 if any check fails.

## Installation

```bash
go build -o curlsc cmd/curlsc/main.go
```

## Usage

### Simple Check
```bash
./curlsc https://google.com https://github.com
```

### Check from File
```bash
./curlsc --file urls.txt --concurrency 50
```

### Piped Input
```bash
cat urls.txt | ./curlsc --format json
```

### Options

| Flag | Description | Default |
|------|-------------|---------|
| `-f, --file` | Path to file with URLs | - |
| `-c, --concurrency` | Max parallel requests | 20 |
| `-t, --timeout` | Per-request timeout | 10s |
| `-m, --method` | HTTP method (GET/HEAD) | GET |
| `--format` | Output format (table/json/csv) | table |
| `-o, --output` | Write results to file | stdout |
| `-k, --insecure` | Skip TLS verification | false |

## License
MIT
