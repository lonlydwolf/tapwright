# Tapwright

A Go CLI for managing a local personal Homebrew tap.

Tapwright is in early development. Its first implemented capability is **offline GitHub repository inspection**. Homebrew tap and formula management are planned but are not implemented yet.

## Repository inspection

The currently supported command is:

```text
tapwright inspect <repository-url>
```

It validates a GitHub repository URL and extracts the owner, repository name, and canonical URL. Inspection makes **no network requests**: it does not check whether a repository exists, is accessible, has releases, or can be installed.

For example, after building Tapwright:

```sh
./tapwright inspect https://github.com/lonlydwolf/tapwright
```

Output (one JSON object followed by a newline):

```json
{
  "owner": "lonlydwolf",
  "repository": "tapwright",
  "canonical_url": "https://github.com/lonlydwolf/tapwright"
}
```

### Accepted URLs

The accepted URL form is exactly `https://github.com/<owner>/<repository>`.

- The scheme must be spelled **`https`** and the host **`github.com`**, both in lowercase.
- No port or embedded credentials are permitted.
- The owner and repository must each be nonempty and contain only ASCII letters (`A–Z`, `a–z`), digits (`0–9`), hyphens (`-`), underscores (`_`), or periods (`.`).
- Neither path component may consist solely of `.` or `..`.
- Path components preserve their original spelling and case. A `.git` suffix is treated as part of the repository name, not removed.
- Percent-encoded path components are not supported.
- Queries and fragments are rejected, including an empty `?` or `#` delimiter.
- Missing or extra path components and trailing slashes are rejected. Other URL forms and malformed URLs are unsupported.

Examples:

| URL                                       | Result                            |
| ----------------------------------------- | --------------------------------- |
| `https://github.com/octocat/Hello-World`  | Accepted                          |
| `https://github.com/octocat/repo_123.git` | Accepted; `.git` is preserved     |
| `HTTPS://github.com/octocat/repo`         | Rejected; scheme spelling differs |
| `https://github.com:443/octocat/repo`     | Rejected; port present            |
| `https://github.com/octocat/%48ello`      | Rejected; encoded path component  |
| `https://github.com/octocat/repo/`        | Rejected; trailing slash          |
| `https://github.com/octocat/repo?`        | Rejected; query delimiter         |

### Output and exit status

- **Success:** Exit status `0`; a single JSON object containing exactly the string fields `owner`, `repository`, and `canonical_url` is written to stdout, followed by a newline. Stderr is empty. JSON field order and whitespace are not part of the contract.
- **Failure:** Exit status `1`; a useful error describing the argument or URL problem is written to stderr. Stdout is empty.

The CLI accepts exactly the `inspect` subcommand with one URL argument. Missing arguments, unknown subcommands, and extra arguments are errors. Flags and help commands are not part of the current interface.

## Getting started

### Requirements

- Go **1.27.2 or newer** (the version declared in `go.mod` is `1.27.2`).

### Build

```sh
git clone https://github.com/lonlydwolf/tapwright.git
cd tapwright
go build -o tapwright .
```

The generated `tapwright` executable is a local build artifact, not source code.

### Run

```sh
./tapwright inspect https://github.com/octocat/Hello-World
```

## Development checks

Run from the repository root:

```sh
go build ./...
go test -v ./...
sh scripts/check-cli.sh
```

The CLI script builds an executable and verifies process exit codes, stdout, and stderr for both successful and rejected invocations. GitHub Actions runs Go build, unit tests, and CLI checks for pull requests targeting `main`.

## Roadmap

Planned future work includes:

- GitHub release discovery and metadata retrieval.
- Homebrew formula generation and maintenance.
- Local personal tap maintenance.

These capabilities are **not yet implemented**.
