# Tapwright

A Go CLI for managing a local personal Homebrew tap.

Tapwright is in early development. It currently supports two commands:

- `tapwright inspect <repository-url>`: **offline GitHub repository inspection**.
- `tapwright create <user> <name>`: **scaffolding a new local tap** under an existing taps directory.

Formula and cask management are planned but are not implemented yet.

## Repository inspection

The command is:

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
{"owner":"lonlydwolf","repository":"tapwright","canonical_url":"https://github.com/lonlydwolf/tapwright"}
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

| URL | Result |
| --- | --- |
| `https://github.com/octocat/Hello-World` | Accepted |
| `https://github.com/octocat/repo_123.git` | Accepted; `.git` is preserved |
| `HTTPS://github.com/octocat/repo` | Rejected; scheme spelling differs |
| `https://github.com:443/octocat/repo` | Rejected; port present |
| `https://github.com/octocat/%48ello` | Rejected; encoded path component |
| `https://github.com/octocat/repo/` | Rejected; trailing slash |
| `https://github.com/octocat/repo?` | Rejected; query delimiter |

### Output and exit status

- **Success:** Exit status `0`; a single JSON object containing exactly the string fields `owner`, `repository`, and `canonical_url` is written to stdout, followed by a newline. Stderr is empty. JSON field order and whitespace are not part of the contract.
- **Failure:** Exit status `1`; a useful error describing the argument or URL problem is written to stderr. Stdout is empty.

`inspect` takes exactly one URL argument. Missing or extra arguments are errors.

## Tap creation

```text
tapwright create <user> <name>
```

It creates a new tap directory at `<taps-root>/<user>/homebrew-<name>`, the layout Homebrew uses for a tap named `<user>/<name>`. You give the short name; Tapwright adds the `homebrew-` prefix.

For example:

```sh
./tapwright create alice tools
```

Output when using the default taps root:

```text
/opt/homebrew/Library/Taps/alice/homebrew-tools
```

### Taps root

- If the `TAPWRIGHT_TAPS_ROOT` environment variable is nonempty, its value is used as the taps root.
- If the variable is unset or empty, the default is `/opt/homebrew/Library/Taps`.
- The taps root must already exist and be an accessible directory.
- Tapwright does not create the taps root automatically.
- If the taps root cannot be opened, `create` fails with a `taps root unavailable` error.

Tapwright does not invoke Homebrew to discover the taps directory. For non-default Homebrew installations or standalone local use, specify the taps root through `TAPWRIGHT_TAPS_ROOT`.

### Names

- `<user>` and `<name>` must each be one or more segments of lowercase ASCII letters (`a–z`) and digits (`0–9`), separated by single hyphens. Leading, trailing, or consecutive hyphens, uppercase letters, underscores, dots, slashes, spaces, and other characters are rejected.
- `<name>` must not start with `homebrew-`, because the prefix is added for you. For example, `tapwright create alice homebrew-tools` is rejected.

### What gets created

```text
<taps-root>/<user>/homebrew-<name>/
├── Casks/
└── .tapwright
```

- The `<user>` directory is created if it does not exist.
- `.tapwright` is a marker file that identifies the directory as a Tapwright-managed tap.
- If the target tap already exists, `create` fails with `tap already exists` and leaves it untouched.
- Filesystem operations are confined to the taps root. A `<user>` directory that is a symlink pointing outside the root is rejected.
- If creating the tap contents fails partway, Tapwright attempts to remove the newly created tap directory. A newly created user directory may remain.
- Tapwright does not initialize a Git repository.
- `create` does not run `brew tap`, and it never modifies an existing tap.

### Output and exit status

- **Success:** Exit status `0`; the path of the new tap directory is written to stdout, followed by a newline. Stderr is empty.
- **Failure:** Exit status `1`; a useful error describing the problem is written to stderr. Stdout is empty.

The output path is constructed using `filepath.Join` and may be relative
or absolute depending on the supplied taps root. It is lexically cleaned
without resolving symlinks, so roots containing symlinks and `..` may
produce an output path that differs from the actual location.

`create` takes exactly two arguments. Missing or extra arguments print `Usage: tapwright create <user> <name>` and exit with status `1`.

## Command-line interface

The CLI accepts exactly the `inspect` and `create` subcommands. Running `tapwright` with no command and running an unknown command are errors (exit status `1`). Flags and help commands are not part of the current interface.

## Getting started

### Requirements

- Go **1.27.2 or newer** (the version declared in `go.mod` is `1.27.2`).
- An existing taps directory when using `create`.

### Build

```sh
git clone https://github.com/lonlydwolf/tapwright.git
cd tapwright
go build -o tapwright .
```

The generated `tapwright` executable is a local build artifact, not source code.

### Run

Inspect a repository:

```sh
./tapwright inspect https://github.com/octocat/Hello-World
```

Create a tap using the default taps root (this writes to your real Homebrew taps directory):

```sh
./tapwright create alice tools
```

Create a tap using a custom taps root:

```sh
mkdir -p "$PWD/taps"
TAPWRIGHT_TAPS_ROOT="$PWD/taps" ./tapwright create alice tools
```

## Development checks

Run from the repository root:

```sh
go build ./...
go test -v ./...
bash scripts/check-cli.sh
```

The CLI script requires Bash. It builds an executable in a temporary directory and verifies process exit codes, stdout, and stderr for both successful and rejected invocations of `inspect` and `create`.

The script uses temporary taps root paths, so it never writes to your real Homebrew taps directory, and removes its temporary files on exit.

## Roadmap

Planned future work includes:

- GitHub release discovery and metadata retrieval.
- Homebrew formula and cask generation and maintenance.
- Maintenance of existing local taps beyond initial creation.

These capabilities are **not yet implemented**.
