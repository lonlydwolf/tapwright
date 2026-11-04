# Tapwright

A Go CLI for managing a local personal Homebrew tap.

Tapwright aims to simplify the maintenance of personal Homebrew taps, from inspecting GitHub repositories to managing Homebrew formulae and releases.

> **Status:** Early development. Homebrew tap management functionality is not yet implemented.

## Features

### Repository Inspection

Validate GitHub repository URLs and extract structured metadata without making network requests.

```bash
tapwright inspect https://github.com/lonlydwolf/tapwright
```

**Output:**

```json
{
  "owner": "lonlydwolf",
  "repository": "tapwright",
  "canonical_url": "https://github.com/lonlydwolf/tapwright"
}
```

The command enforces strict URL validation:

- Requires HTTPS and the `github.com` host.
- Rejects embedded credentials.
- Rejects queries and fragments, including empty delimiters.
- Requires exactly two non-empty path components (`owner` and `repository`).
- Rejects additional path components and `.` or `..` segments.

**Exit behavior:**

| Exit code | Behavior                                              |
| --------- | ----------------------------------------------------- |
| `0`       | Outputs JSON to `stdout`.                             |
| `1`       | Outputs an error to `stderr`, leaving `stdout` empty. |

## Getting Started

### Requirements

- Go 1.27.2 or newer

### Build

```bash
git clone https://github.com/lonlydwolf/tapwright.git
cd tapwright
go build -o tapwright .
```

### Usage

```bash
./tapwright inspect https://github.com/<owner>/<repository>
```

## Development

Build and verify all packages:

```bash
go build ./...
```

Run unit tests:

```bash
go test -v ./...
```

Run CLI integration tests:

```bash
sh scripts/check-cli.sh
```

GitHub Actions automatically runs these checks on pull requests targeting `main`.

## Roadmap

Tapwright is being developed incrementally toward its primary goal of managing a local personal Homebrew tap.

Planned capabilities include:

- GitHub release discovery and metadata retrieval.
- Homebrew formula generation and maintenance.
- Local
